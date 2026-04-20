package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mahimairaja/vbench/internal/adapter"
	"github.com/mahimairaja/vbench/internal/concurrency"
	"github.com/mahimairaja/vbench/internal/schema"
)

// index writes every conversation turn into the provider via add_message, and
// records per-item write latency artifacts. It also issues a reset() per item
// to guarantee isolation between callers (voice-realistic: each caller is a
// separate session with its own memory scope).
func (p *Pipeline) index(ctx context.Context) error {
	dir := stageDir(p.ArtifactDir, "index")
	if stageComplete(dir) {
		if p.Verbose {
			fmt.Fprintln(os.Stderr, "[index] resume: already complete")
		}
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	client := p.Side.Client()
	for _, item := range p.Items {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		// Ensure a clean user namespace before writing.
		userID := firstUserID(item)
		if err := client.Reset(ctx, adapter.ResetRequest{UserID: userID}); err != nil {
			return fmt.Errorf("reset %s: %w", userID, err)
		}

		var latencies []float64
		var results []map[string]interface{}
		for _, turn := range item.Conversation {
			start := time.Now()
			wr, err := client.AddMessage(ctx, adapter.AddMessageRequest{
				UserID:    turn.UserID,
				SessionID: turn.SessionID,
				Role:      turn.Role,
				Content:   turn.Content,
			})
			elapsedMs := float64(time.Since(start).Microseconds()) / 1000.0
			if err != nil {
				return fmt.Errorf("add_message turn=%s: %w", turn.TurnID, err)
			}
			latencies = append(latencies, elapsedMs)
			results = append(results, map[string]interface{}{
				"turn_id":     turn.TurnID,
				"provider_id": wr.ProviderID,
				"latency_ms":  elapsedMs,
			})
		}
		artifact := schema.IndexArtifact{
			ItemID:         item.ItemID,
			Provider:       p.Config.Provider.Name,
			WriteResults:   results,
			TotalLatencyMs: sum(latencies),
			P50LatencyMs:   concurrency.Percentile(latencies, 50),
			P95LatencyMs:   concurrency.Percentile(latencies, 95),
		}
		if err := writeJSON(filepath.Join(dir, item.ItemID+".json"), artifact); err != nil {
			return err
		}
		if p.Verbose {
			fmt.Fprintf(os.Stderr, "[index] %s: %d turns, p95 %.1f ms\n", item.ItemID, len(latencies), artifact.P95LatencyMs)
		}
	}
	return markComplete(dir)
}

func sum(xs []float64) float64 {
	var total float64
	for _, v := range xs {
		total += v
	}
	return total
}

func firstUserID(item schema.BenchmarkItem) string {
	for _, t := range item.Conversation {
		if t.UserID != "" {
			return t.UserID
		}
	}
	return "user_" + item.ItemID
}

// lastSessionID returns the most recent session_id in the conversation. Search
// happens at the end of the conversation, so scoping reads to the latest
// session is the voice-realistic default.
func lastSessionID(item schema.BenchmarkItem) string {
	for i := len(item.Conversation) - 1; i >= 0; i-- {
		if item.Conversation[i].SessionID != "" {
			return item.Conversation[i].SessionID
		}
	}
	return "session_" + item.ItemID
}
