package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// ingest writes the normalised benchmark items to disk so downstream stages
// have a stable input artifact. For LoCoMo (already in-memory) this is mostly
// a consistency checkpoint, but it keeps all stages uniform and resumable.
func (p *Pipeline) ingest(ctx context.Context) error {
	dir := stageDir(p.ArtifactDir, "ingest")
	if stageComplete(dir) {
		if p.Verbose {
			fmt.Fprintln(os.Stderr, "[ingest] resume: already complete")
		}
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, item := range p.Items {
		if err := writeJSON(filepath.Join(dir, item.ItemID+".json"), item); err != nil {
			return err
		}
	}
	if p.Verbose {
		fmt.Fprintf(os.Stderr, "[ingest] wrote %d items\n", len(p.Items))
	}
	return markComplete(dir)
}
