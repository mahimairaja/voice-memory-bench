package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mahimairaja/vbench/internal/concurrency"
	"github.com/mahimairaja/vbench/internal/schema"
)

// evaluate joins search artifacts (for latency) with answer artifacts (for
// text) and runs the judge LLM. It emits per-item evaluation JSON and returns
// the aggregate MemScore for this concurrency level.
func (p *Pipeline) evaluate(ctx context.Context, level int) (schema.MemScore, error) {
	outDir := filepath.Join(stageDir(p.ArtifactDir, "evaluate"), fmt.Sprintf("%dx", level))
	searchDir := filepath.Join(stageDir(p.ArtifactDir, "search"), fmt.Sprintf("%dx", level))
	answerDir := filepath.Join(stageDir(p.ArtifactDir, "answer"), fmt.Sprintf("%dx", level))
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return schema.MemScore{}, err
	}

	qMap := buildQuestionMap(p.Items)

	type searchRec struct {
		ItemID         string  `json:"item_id"`
		QuestionID     string  `json:"question_id"`
		LatencyMs      float64 `json:"latency_ms"`
		TokenFootprint int     `json:"token_footprint"`
	}
	type answerRec struct {
		ItemID     string `json:"item_id"`
		QuestionID string `json:"question_id"`
		Completion string `json:"completion"`
	}

	// Collect per-question latency + token footprint + judged score.
	var latencies []float64
	var tokens []int
	perItem := make(map[string]*schema.EvaluationArtifact)

	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return schema.MemScore{}, err
	}
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == ".complete" || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		if ctx.Err() != nil {
			return schema.MemScore{}, ctx.Err()
		}

		sRaw, err := os.ReadFile(filepath.Join(searchDir, entry.Name()))
		if err != nil {
			return schema.MemScore{}, err
		}
		var s searchRec
		if err := json.Unmarshal(sRaw, &s); err != nil {
			return schema.MemScore{}, fmt.Errorf("decode search artifact: %w", err)
		}
		aRaw, err := os.ReadFile(filepath.Join(answerDir, entry.Name()))
		if err != nil {
			return schema.MemScore{}, err
		}
		var a answerRec
		if err := json.Unmarshal(aRaw, &a); err != nil {
			return schema.MemScore{}, fmt.Errorf("decode answer artifact: %w", err)
		}
		q := qMap[s.QuestionID]
		verdict, err := p.JudgeLLM.Judge(ctx, q.Question, q.ReferenceAnswer, a.Completion)
		if err != nil {
			return schema.MemScore{}, fmt.Errorf("judge %s: %w", s.QuestionID, err)
		}

		latencies = append(latencies, s.LatencyMs)
		tokens = append(tokens, s.TokenFootprint)

		ev, ok := perItem[s.ItemID]
		if !ok {
			ev = &schema.EvaluationArtifact{
				ItemID:   s.ItemID,
				Provider: p.Config.Provider.Name,
				Dataset:  p.Config.Dataset.Name,
			}
			perItem[s.ItemID] = ev
		}
		ev.PerQuestionScores = append(ev.PerQuestionScores, schema.QuestionScore{
			QuestionID: s.QuestionID,
			Score:      verdict.Score,
			Rationale:  verdict.Rationale,
		})
	}

	// Aggregate.
	var totalScore float64
	var totalCount int
	for _, ev := range perItem {
		for _, qs := range ev.PerQuestionScores {
			totalScore += qs.Score
			totalCount++
		}
	}
	avgQuality := 0.0
	if totalCount > 0 {
		avgQuality = totalScore / float64(totalCount)
	}
	ms := schema.MemScore{
		Quality:           avgQuality,
		LatencyP50Ms:      concurrency.Percentile(latencies, 50),
		LatencyP95Ms:      concurrency.Percentile(latencies, 95),
		LatencyP99Ms:      concurrency.Percentile(latencies, 99),
		CostPerItem:       0,
		TokenFootprintP50: medianInt(tokens),
		Concurrency:       level,
		NumQuestions:      totalCount,
	}

	// Per-item evaluation artifacts + run-level score
	for id, ev := range perItem {
		ev.MemScore = ms
		if err := writeJSON(filepath.Join(outDir, id+".json"), ev); err != nil {
			return schema.MemScore{}, err
		}
	}
	if err := writeJSON(filepath.Join(outDir, "_memscore.json"), ms); err != nil {
		return schema.MemScore{}, err
	}
	return ms, markComplete(outDir)
}

func medianInt(xs []int) int {
	if len(xs) == 0 {
		return 0
	}
	as := make([]float64, len(xs))
	for i, v := range xs {
		as[i] = float64(v)
	}
	return int(concurrency.Percentile(as, 50))
}
