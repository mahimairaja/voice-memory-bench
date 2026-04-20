package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/mahimairaja/vbench/internal/schema"
)

// Client wraps the OpenAI Chat Completions endpoint. It is sufficient for both
// the answer LLM and the judge LLM roles in the MVP; only the model/params
// differ between roles.
type Client struct {
	cfg     schema.LLMConfig
	apiKey  string
	baseURL string
	http    *http.Client
}

// New constructs a Client from config + API key. apiKey may be empty when the
// caller is pointing at a local endpoint that does not require auth.
func New(cfg schema.LLMConfig, apiKey string) *Client {
	base := cfg.BaseURL
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	return &Client{
		cfg:     cfg,
		apiKey:  apiKey,
		baseURL: base,
		http: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ChatMessage is one element of a chat completion request.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Completion is the stripped-down response we carry through the pipeline.
type Completion struct {
	Text             string
	PromptTokens     int
	CompletionTokens int
}

type chatReq struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Seed        *int          `json:"seed,omitempty"`
}

type chatResp struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage"`
}

// Complete runs a chat completion.
func (c *Client) Complete(ctx context.Context, messages []ChatMessage) (*Completion, error) {
	body := chatReq{
		Model:       c.cfg.Model,
		Messages:    messages,
		Temperature: c.cfg.Temperature,
		MaxTokens:   c.cfg.MaxTokens,
	}
	if c.cfg.Seed != 0 {
		s := c.cfg.Seed
		body.Seed = &s
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai %s: %d %s", c.cfg.Model, resp.StatusCode, string(b))
	}
	var parsed chatResp
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode openai response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}
	return &Completion{
		Text:             parsed.Choices[0].Message.Content,
		PromptTokens:     parsed.Usage.PromptTokens,
		CompletionTokens: parsed.Usage.CompletionTokens,
	}, nil
}

// Answer composes the answer-stage prompt from injected memory + question.
func (c *Client) Answer(ctx context.Context, memoryPayload, question string) (*Completion, error) {
	system := "You are a helpful voice assistant. Answer the user's question concisely using only the information in the provided memory context. If the memory does not contain the answer, say you don't know."
	user := fmt.Sprintf("Memory context:\n%s\n\nQuestion: %s\nAnswer:", memoryPayload, question)
	return c.Complete(ctx, []ChatMessage{
		{Role: "system", Content: system},
		{Role: "user", Content: user},
	})
}

// JudgeVerdict is the structured output we parse out of the judge LLM.
type JudgeVerdict struct {
	Score     float64
	Rationale string
}

var judgeScorePattern = regexp.MustCompile(`(?i)score\s*[:=]\s*([0-1](?:\.\d+)?)`)

// Judge runs the judge LLM on (question, reference answer, candidate answer)
// and returns a [0,1] score plus the judge's rationale text. A malformed
// response (no parseable Score line) is surfaced as an error rather than
// silently returning 0.0, which would corrupt the aggregate quality.
func (c *Client) Judge(ctx context.Context, question, reference, candidate string) (*JudgeVerdict, error) {
	system := "You grade answers for factual correctness against a reference. Respond with exactly two lines:\nScore: <number between 0 and 1, where 1 means the candidate matches the reference and 0 means it does not>\nRationale: <one sentence>"
	user := fmt.Sprintf("Question: %s\nReference answer: %s\nCandidate answer: %s", question, reference, candidate)
	comp, err := c.Complete(ctx, []ChatMessage{
		{Role: "system", Content: system},
		{Role: "user", Content: user},
	})
	if err != nil {
		return nil, err
	}
	m := judgeScorePattern.FindStringSubmatch(comp.Text)
	if len(m) != 2 {
		return nil, fmt.Errorf("judge response has no parseable Score line: %q", comp.Text)
	}
	f, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return nil, fmt.Errorf("judge score %q is not a number: %w", m[1], err)
	}
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}
	return &JudgeVerdict{Score: f, Rationale: comp.Text}, nil
}
