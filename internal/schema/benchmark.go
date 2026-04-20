package schema

import "time"

// ConversationTurn is a single turn in a multi-turn conversation.
type ConversationTurn struct {
	TurnID    string            `json:"turn_id"`
	SessionID string            `json:"session_id"`
	UserID    string            `json:"user_id"`
	Role      string            `json:"role"`
	Content   string            `json:"content"`
	Timestamp *time.Time        `json:"timestamp,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// EvaluationQuestion is a question derived from the conversation with a reference answer.
type EvaluationQuestion struct {
	QuestionID      string `json:"question_id"`
	Question        string `json:"question"`
	ReferenceAnswer string `json:"reference_answer"`
	QuestionType    string `json:"question_type,omitempty"`
}

// BenchmarkItem is one unit of work: a conversation plus evaluation questions.
type BenchmarkItem struct {
	ItemID       string               `json:"item_id"`
	Dataset      string               `json:"dataset"`
	Subset       string               `json:"subset,omitempty"`
	Conversation []ConversationTurn   `json:"conversation"`
	Questions    []EvaluationQuestion `json:"questions"`
	Metadata     map[string]string    `json:"metadata,omitempty"`
}
