package dataset

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mahimairaja/vbench/internal/schema"
)

const (
	locomoURL      = "https://huggingface.co/datasets/snap-research/locomo/resolve/main/locomo10_test.json"
	locomoFileName = "locomo10_test.json"
	locomoDir      = "locomo"
)

// LoCoMo is the loader for the LoCoMo (Long-Horizon Conversational Memory) dataset.
// Source: https://huggingface.co/datasets/snap-research/locomo
// Paper: Maharana et al., ACL 2024. https://arxiv.org/abs/2402.17753
type LoCoMo struct{}

func (l *LoCoMo) Name() string { return "locomo" }

// dataPath returns the local JSON path.
func (l *LoCoMo) dataPath(cacheDir string) string {
	return filepath.Join(cacheDir, locomoDir, locomoFileName)
}

// IsCached reports whether the LoCoMo JSON exists and is non-empty.
func (l *LoCoMo) IsCached(cacheDir string) bool {
	info, err := os.Stat(l.dataPath(cacheDir))
	return err == nil && info.Size() > 0
}

// Download fetches the LoCoMo JSON into cacheDir. Cancellable via ctx.
// Surfaces Close() errors so a failed flush is not silently accepted into the
// cache — the next `IsCached` call could otherwise see a truncated file and
// skip redownloading.
func (l *LoCoMo) Download(ctx context.Context, cacheDir string) (err error) {
	dst := l.dataPath(cacheDir)
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir cache: %w", err)
	}
	if l.IsCached(cacheDir) {
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, locomoURL, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download locomo: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download locomo: HTTP %d", resp.StatusCode)
	}
	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create %s: %w", dst, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close %s: %w", dst, cerr)
		}
	}()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("write %s: %w", dst, err)
	}
	return nil
}

// Raw upstream structure. LoCoMo stores a list of samples; each sample has
// multi-session dialogue plus a list of QA pairs. The schema on HuggingFace
// uses "conversation" with numbered session keys and "qa" for questions.
type locomoRaw struct {
	SampleID     string          `json:"sample_id"`
	Conversation json.RawMessage `json:"conversation"`
	QA           []locomoQARaw   `json:"qa"`
}

type locomoQARaw struct {
	Question string      `json:"question"`
	Answer   interface{} `json:"answer"`
	Category interface{} `json:"category"`
}

// A session turn as it appears under conversation.session_N.
type locomoTurnRaw struct {
	DiaID   string `json:"dia_id"`
	Speaker string `json:"speaker"`
	Text    string `json:"text"`
}

// Load parses the LoCoMo JSON into a slice of BenchmarkItems.
// Each sample becomes one BenchmarkItem; sessions are flattened into a single
// conversation list with session_id preserved on each turn.
func (l *LoCoMo) Load(cacheDir, subset string, maxItems int) (items []schema.BenchmarkItem, err error) {
	if !l.IsCached(cacheDir) {
		return nil, fmt.Errorf("locomo not cached at %s; run `vbench datasets download locomo` first", l.dataPath(cacheDir))
	}
	f, err := os.Open(l.dataPath(cacheDir))
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	var raw []locomoRaw
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode locomo json: %w", err)
	}

	for i, r := range raw {
		if maxItems > 0 && len(items) >= maxItems {
			break
		}
		itemID := SafeID(r.SampleID)
		if itemID == "_" || itemID == "" {
			itemID = fmt.Sprintf("locomo_%d", i)
		}
		userID := "locomo_user_" + itemID

		conversation, err := flattenSessions(r.Conversation, itemID, userID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", itemID, err)
		}

		questions := make([]schema.EvaluationQuestion, 0, len(r.QA))
		for qi, qa := range r.QA {
			ans := stringifyAnswer(qa.Answer)
			cat := ""
			if s, ok := qa.Category.(string); ok {
				cat = s
			} else if n, ok := qa.Category.(float64); ok {
				cat = strconv.FormatFloat(n, 'f', -1, 64)
			}
			questions = append(questions, schema.EvaluationQuestion{
				QuestionID:      SafeID(fmt.Sprintf("%s_q%d", itemID, qi)),
				Question:        qa.Question,
				ReferenceAnswer: ans,
				QuestionType:    cat,
			})
		}

		items = append(items, schema.BenchmarkItem{
			ItemID:       itemID,
			Dataset:      "locomo",
			Subset:       subset,
			Conversation: conversation,
			Questions:    questions,
		})
	}
	return items, nil
}

// flattenSessions handles both historical shapes of the conversation field:
//  1. A map with keys like "session_1", "session_1_date_time", "speaker_a", ...
//     where session_N is a list of turns.
//  2. An array of sessions.
func flattenSessions(raw json.RawMessage, itemID, userID string) ([]schema.ConversationTurn, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	// Shape 1: object with session_N keys.
	var asMap map[string]json.RawMessage
	if err := json.Unmarshal(raw, &asMap); err == nil {
		return flattenSessionMap(asMap, itemID, userID)
	}
	// Shape 2: array of sessions, each a list of turns.
	var asArray []json.RawMessage
	if err := json.Unmarshal(raw, &asArray); err == nil {
		var turns []schema.ConversationTurn
		for i, s := range asArray {
			sessionID := fmt.Sprintf("%s_session_%d", itemID, i+1)
			ts, err := decodeTurns(s, sessionID, userID)
			if err != nil {
				return nil, err
			}
			turns = append(turns, ts...)
		}
		return turns, nil
	}
	return nil, fmt.Errorf("conversation has unexpected shape")
}

func flattenSessionMap(m map[string]json.RawMessage, itemID, userID string) ([]schema.ConversationTurn, error) {
	// Preserve natural session order: session_1, session_2, ...
	var out []schema.ConversationTurn
	for i := 1; ; i++ {
		key := fmt.Sprintf("session_%d", i)
		v, ok := m[key]
		if !ok {
			break
		}
		sessionID := fmt.Sprintf("%s_%s", itemID, key)
		ts, err := decodeTurns(v, sessionID, userID)
		if err != nil {
			return nil, err
		}
		out = append(out, ts...)
	}
	return out, nil
}

func decodeTurns(raw json.RawMessage, sessionID, userID string) ([]schema.ConversationTurn, error) {
	var rawTurns []locomoTurnRaw
	if err := json.Unmarshal(raw, &rawTurns); err != nil {
		return nil, fmt.Errorf("decode session turns: %w", err)
	}
	out := make([]schema.ConversationTurn, 0, len(rawTurns))
	for i, t := range rawTurns {
		turnID := t.DiaID
		if turnID == "" {
			turnID = fmt.Sprintf("%s_t%d", sessionID, i)
		}
		out = append(out, schema.ConversationTurn{
			TurnID:    turnID,
			SessionID: sessionID,
			UserID:    userID,
			Role:      roleFromSpeaker(t.Speaker, i),
			Content:   t.Text,
			Metadata:  map[string]string{"speaker": t.Speaker},
		})
	}
	return out, nil
}

// roleFromSpeaker maps a LoCoMo speaker label to a voice-agent role. LoCoMo
// dialogues are peer-to-peer (speaker_a / speaker_b or similar); we treat the
// "a" side as user and the "b" side as assistant. When the speaker label is
// missing, we fall back to alternating by turn index so the write path still
// produces a coherent role sequence.
func roleFromSpeaker(speaker string, turnIdx int) string {
	s := strings.ToLower(strings.TrimSpace(speaker))
	switch {
	case strings.HasSuffix(s, "_b"), s == "b":
		return "assistant"
	case strings.HasSuffix(s, "_a"), s == "a":
		return "user"
	case s == "":
		if turnIdx%2 == 1 {
			return "assistant"
		}
		return "user"
	default:
		return "user"
	}
}

func stringifyAnswer(a interface{}) string {
	switch v := a.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
}
