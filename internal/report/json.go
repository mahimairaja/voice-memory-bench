package report

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mahimairaja/vbench/internal/schema"
)

// Report is the top-level structure written to memscore.json.
type Report struct {
	RunID    string             `json:"run_id"`
	Provider string             `json:"provider"`
	Dataset  string             `json:"dataset"`
	Levels   []LevelReport      `json:"levels"`
	Headline []string           `json:"headline"`
	MemScore []schema.MemScore  `json:"memscore"`
}

// LevelReport bundles one concurrency level's MemScore with its verdict.
type LevelReport struct {
	Concurrency int             `json:"concurrency"`
	Verdict     Verdict         `json:"verdict"`
	MemScore    schema.MemScore `json:"memscore"`
}

// WriteJSON emits the MemScore report for all concurrency levels. Surfaces
// Close() errors after a successful encode so a failed flush on the final
// report is not silently swallowed.
func WriteJSON(path, provider, dataset, runID string, levels []schema.MemScore) (err error) {
	rep := Report{
		RunID:    runID,
		Provider: provider,
		Dataset:  dataset,
		MemScore: levels,
	}
	for _, ms := range levels {
		rep.Levels = append(rep.Levels, LevelReport{
			Concurrency: ms.Concurrency,
			Verdict:     ClassifyP95(ms.LatencyP95Ms),
			MemScore:    ms,
		})
		rep.Headline = append(rep.Headline, Headline(provider, ms))
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}
