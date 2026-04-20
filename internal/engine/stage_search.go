package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mahimairaja/vbench/internal/adapter"
	"github.com/mahimairaja/vbench/internal/concurrency"
)

// search is the voice-critical read path. We fire one search per (item,
// question) with `concurrency` workers in flight, measuring wall-clock latency
// in the Go engine (authoritative) rather than trusting the sidecar-reported
// number.
//
// Results are bucketed per concurrency level into search/<level>/ so a single
// run can compare p95 at 1x vs 4x without re-indexing.
func (p *Pipeline) search(ctx context.Context, level int) error {
	dir := filepath.Join(stageDir(p.ArtifactDir, "search"), fmt.Sprintf("%dx", level))
	if stageComplete(dir) {
		if p.Verbose {
			fmt.Fprintf(os.Stderr, "[search @%dx] resume: already complete\n", level)
		}
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	type searchJob struct {
		itemID     string
		userID     string
		sessionID  string
		questionID string
		query      string
	}
	var jobs []searchJob
	for _, item := range p.Items {
		userID := firstUserID(item)
		sessionID := lastSessionID(item)
		for _, q := range item.Questions {
			jobs = append(jobs, searchJob{
				itemID:     item.ItemID,
				userID:     userID,
				sessionID:  sessionID,
				questionID: q.QuestionID,
				query:      q.Question,
			})
		}
	}

	client := p.Side.Client()
	results := concurrency.Run(ctx, len(jobs), level, func(ctx context.Context, i int) (searchOutput, error) {
		j := jobs[i]
		start := time.Now()
		res, err := client.Search(ctx, adapter.SearchRequest{
			UserID:    j.userID,
			SessionID: j.sessionID,
			Query:     j.query,
			Mode:      adapter.ModeSemantic,
			TopK:      10,
		})
		elapsedMs := float64(time.Since(start).Microseconds()) / 1000.0
		if err != nil {
			return searchOutput{}, err
		}
		payload := renderMemoryPayload(res.Items)
		return searchOutput{
			ItemID:         j.itemID,
			QuestionID:     j.questionID,
			Provider:       p.Config.Provider.Name,
			Concurrency:    level,
			LatencyMs:      elapsedMs,
			Result:         res,
			MemoryPayload:  payload,
			TokenFootprint: approxTokenCount(payload),
		}, nil
	})

	for _, r := range results {
		if r.Err != nil {
			return fmt.Errorf("search job %d: %w", r.Index, r.Err)
		}
		o := r.Value
		artifact := map[string]interface{}{
			"item_id":          o.ItemID,
			"question_id":      o.QuestionID,
			"provider":         o.Provider,
			"concurrency":      o.Concurrency,
			"latency_ms":       o.LatencyMs,
			"retrieval_result": o.Result,
			"memory_payload":   o.MemoryPayload,
			"token_footprint":  o.TokenFootprint,
		}
		if err := writeJSON(filepath.Join(dir, o.ItemID+"__"+o.QuestionID+".json"), artifact); err != nil {
			return err
		}
	}
	if p.Verbose {
		var lats []float64
		for _, r := range results {
			lats = append(lats, r.Value.LatencyMs)
		}
		fmt.Fprintf(os.Stderr, "[search @%dx] n=%d p50=%.1f p95=%.1f p99=%.1f ms\n",
			level, len(lats),
			concurrency.Percentile(lats, 50),
			concurrency.Percentile(lats, 95),
			concurrency.Percentile(lats, 99),
		)
	}
	return markComplete(dir)
}

type searchOutput struct {
	ItemID         string
	QuestionID     string
	Provider       string
	Concurrency    int
	LatencyMs      float64
	Result         *adapter.RetrievalResult
	MemoryPayload  string
	TokenFootprint int
}

// renderMemoryPayload builds the exact text block we would inject into a voice
// agent's prompt for this retrieval.
func renderMemoryPayload(items []adapter.MemoryItem) string {
	if len(items) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, it := range items {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("- ")
		sb.WriteString(it.Content)
	}
	return sb.String()
}

// approxTokenCount is a whitespace-based token approximation. The MVP does not
// require model-accurate counts to produce a voice verdict.
func approxTokenCount(s string) int {
	if s == "" {
		return 0
	}
	return len(strings.Fields(s))
}

