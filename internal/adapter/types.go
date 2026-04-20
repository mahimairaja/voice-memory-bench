package adapter

// RetrievalMode enumerates the retrieval strategies a provider may support.
type RetrievalMode string

const (
	ModeSemantic RetrievalMode = "semantic"
	ModeKeyword  RetrievalMode = "keyword"
	ModeTemporal RetrievalMode = "temporal"
	ModeGraph    RetrievalMode = "graph"
	ModeHybrid   RetrievalMode = "hybrid"
)

// Capabilities is the sidecar-reported capability descriptor.
type Capabilities struct {
	ProviderName            string                 `json:"provider_name"`
	ProviderVersion         string                 `json:"provider_version"`
	SupportedRetrievalModes []RetrievalMode        `json:"supported_retrieval_modes"`
	BackingStore            string                 `json:"backing_store"`
	SupportsUserScoping     bool                   `json:"supports_user_scoping"`
	SupportsSessionScoping  bool                   `json:"supports_session_scoping"`
	DeclaredCostModel       string                 `json:"declared_cost_model,omitempty"`
	Extra                   map[string]interface{} `json:"extra,omitempty"`
}

// WriteResult is returned by add_message / add_fact.
type WriteResult struct {
	ProviderID    string                 `json:"provider_id,omitempty"`
	LatencyMs     float64                `json:"latency_ms"`
	TokensWritten int                    `json:"tokens_written,omitempty"`
	Extra         map[string]interface{} `json:"extra,omitempty"`
}

// MemoryItem is one result of a retrieval call.
type MemoryItem struct {
	ItemID    string                 `json:"item_id"`
	Content   string                 `json:"content"`
	Score     *float64               `json:"score,omitempty"`
	CreatedAt *string                `json:"created_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// RetrievalResult is the sidecar-reported retrieval output.
// NOTE: LatencyMs is advisory only. The engine measures the Go-side wall-clock
// latency around each HTTP call; that is the authoritative figure.
type RetrievalResult struct {
	Items          []MemoryItem           `json:"items"`
	LatencyMs      float64                `json:"latency_ms"`
	RetrievalMode  RetrievalMode          `json:"retrieval_mode"`
	TokenFootprint int                    `json:"token_footprint,omitempty"`
	Extra          map[string]interface{} `json:"extra,omitempty"`
}

// AddMessageRequest is the request body for POST /add_message.
type AddMessageRequest struct {
	UserID    string                 `json:"user_id"`
	SessionID string                 `json:"session_id"`
	Role      string                 `json:"role"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AddFactRequest is the request body for POST /add_fact.
type AddFactRequest struct {
	UserID    string                 `json:"user_id"`
	SessionID string                 `json:"session_id"`
	Fact      string                 `json:"fact"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// SearchRequest is the request body for POST /search.
type SearchRequest struct {
	UserID    string                 `json:"user_id"`
	SessionID string                 `json:"session_id"`
	Query     string                 `json:"query"`
	Mode      RetrievalMode          `json:"mode"`
	TopK      int                    `json:"top_k"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
}

// EnumerateRequest is the request body for POST /enumerate.
type EnumerateRequest struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id,omitempty"`
}

// ResetRequest is the request body for POST /reset.
type ResetRequest struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id,omitempty"`
}
