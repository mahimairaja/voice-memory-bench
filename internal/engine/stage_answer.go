package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mahimairaja/vbench/internal/schema"
)

// answer reads the search artifacts for this concurrency level and runs the
// answer LLM on each (memory_payload, question) pair. Latency here is
// informational; the voice verdict is driven by the search-stage p95.
func (p *Pipeline) answer(ctx context.Context, level int) error {
	dir := filepath.Join(stageDir(p.ArtifactDir, "answer"), fmt.Sprintf("%dx", level))
	searchDir := filepath.Join(stageDir(p.ArtifactDir, "search"), fmt.Sprintf("%dx", level))
	if stageComplete(dir) {
		if p.Verbose {
			fmt.Fprintf(os.Stderr, "[answer @%dx] resume: already complete\n", level)
		}
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	qMap := buildQuestionMap(p.Items)

	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return fmt.Errorf("read search dir: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == ".complete" || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		data, err := os.ReadFile(filepath.Join(searchDir, entry.Name()))
		if err != nil {
			return err
		}
		var sa struct {
			ItemID        string `json:"item_id"`
			QuestionID    string `json:"question_id"`
			MemoryPayload string `json:"memory_payload"`
		}
		if err := json.Unmarshal(data, &sa); err != nil {
			return fmt.Errorf("decode search artifact %s: %w", entry.Name(), err)
		}
		q, ok := qMap[sa.QuestionID]
		if !ok {
			return fmt.Errorf("question %s not found in dataset", sa.QuestionID)
		}
		start := time.Now()
		comp, err := p.AnswerLLM.Answer(ctx, sa.MemoryPayload, q.Question)
		elapsedMs := float64(time.Since(start).Microseconds()) / 1000.0
		if err != nil {
			return fmt.Errorf("answer LLM for %s: %w", sa.QuestionID, err)
		}
		artifact := schema.AnswerArtifact{
			ItemID:           sa.ItemID,
			QuestionID:       sa.QuestionID,
			Provider:         p.Config.Provider.Name,
			Prompt:           q.Question,
			Completion:       comp.Text,
			PromptTokens:     comp.PromptTokens,
			CompletionTokens: comp.CompletionTokens,
			LatencyMs:        elapsedMs,
		}
		if err := writeJSON(filepath.Join(dir, entry.Name()), artifact); err != nil {
			return err
		}
	}
	return markComplete(dir)
}

func buildQuestionMap(items []schema.BenchmarkItem) map[string]schema.EvaluationQuestion {
	out := make(map[string]schema.EvaluationQuestion)
	for _, item := range items {
		for _, q := range item.Questions {
			out[q.QuestionID] = q
		}
	}
	return out
}
