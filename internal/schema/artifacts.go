package schema

// IndexArtifact is written per-item by the index stage.
type IndexArtifact struct {
	ItemID         string                   `json:"item_id"`
	Provider       string                   `json:"provider"`
	WriteResults   []map[string]interface{} `json:"write_results"`
	TotalLatencyMs float64                  `json:"total_latency_ms"`
	P50LatencyMs   float64                  `json:"p50_latency_ms"`
	P95LatencyMs   float64                  `json:"p95_latency_ms"`
}

// SearchArtifact is written per-question by the search stage.
// memory_payload is the exact text that would be injected into the prompt.
type SearchArtifact struct {
	ItemID          string                 `json:"item_id"`
	QuestionID      string                 `json:"question_id"`
	Provider        string                 `json:"provider"`
	ConcurrencyLvl  int                    `json:"concurrency"`
	LatencyMs       float64                `json:"latency_ms"`
	RetrievalResult map[string]interface{} `json:"retrieval_result"`
	MemoryPayload   string                 `json:"memory_payload"`
	TokenFootprint  int                    `json:"token_footprint"`
}

// AnswerArtifact is written per-question by the answer stage.
type AnswerArtifact struct {
	ItemID           string  `json:"item_id"`
	QuestionID       string  `json:"question_id"`
	Provider         string  `json:"provider"`
	Prompt           string  `json:"prompt"`
	Completion       string  `json:"completion"`
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	LatencyMs        float64 `json:"latency_ms"`
}

// QuestionScore is the per-question judge verdict.
type QuestionScore struct {
	QuestionID string  `json:"question_id"`
	Score      float64 `json:"score"`
	Rationale  string  `json:"rationale,omitempty"`
}

// EvaluationArtifact is written per-item by the evaluate stage.
type EvaluationArtifact struct {
	ItemID            string          `json:"item_id"`
	Provider          string          `json:"provider"`
	Dataset           string          `json:"dataset"`
	PerQuestionScores []QuestionScore `json:"per_question_scores"`
	MemScore          MemScore        `json:"memscore"`
}
