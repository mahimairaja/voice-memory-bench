package report

import (
	"fmt"
	"strings"

	"github.com/mahimairaja/vbench/internal/schema"
)

// Verdict is the three-level voice-fitness classification.
type Verdict string

const (
	VerdictExcellent  Verdict = "EXCELLENT"
	VerdictAcceptable Verdict = "ACCEPTABLE"
	VerdictFail       Verdict = "FAIL"
)

// ClassifyP95 maps a search-stage p95 latency (ms) to a voice verdict.
// Thresholds confirmed with the user:
//   - p95 < 300 ms           → EXCELLENT
//   - 300 <= p95 <= 500 ms   → ACCEPTABLE
//   - p95 > 500 ms           → FAIL
func ClassifyP95(p95Ms float64) Verdict {
	switch {
	case p95Ms < 300:
		return VerdictExcellent
	case p95Ms <= 500:
		return VerdictAcceptable
	default:
		return VerdictFail
	}
}

// Headline renders the one-line voice verdict.
//
//	mem0 @ 4x: ACCEPTABLE (p95 = 380 ms, quality = 0.72, tokens = 420)
func Headline(provider string, ms schema.MemScore) string {
	return fmt.Sprintf("%s @ %dx: %s (p95 = %.0f ms, quality = %.2f, tokens = %d)",
		provider,
		ms.Concurrency,
		ClassifyP95(ms.LatencyP95Ms),
		ms.LatencyP95Ms,
		ms.Quality,
		ms.TokenFootprintP50,
	)
}

// Summary renders one line per concurrency level, joined by newlines.
func Summary(provider string, levels []schema.MemScore) string {
	var sb strings.Builder
	for i, ms := range levels {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(Headline(provider, ms))
	}
	return sb.String()
}
