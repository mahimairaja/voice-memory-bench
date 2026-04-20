package schema

// MemScore is the quality/latency/cost triple per (provider, concurrency) level.
// Never collapsed to a scalar — the three axes are reported side by side.
type MemScore struct {
	Quality            float64 `json:"quality"`
	LatencyP50Ms       float64 `json:"latency_p50_ms"`
	LatencyP95Ms       float64 `json:"latency_p95_ms"`
	LatencyP99Ms       float64 `json:"latency_p99_ms"`
	CostPerItem        float64 `json:"cost_per_item"`
	TokenFootprintP50  int     `json:"token_footprint_p50"`
	Concurrency        int     `json:"concurrency"`
	NumQuestions       int     `json:"num_questions"`
}
