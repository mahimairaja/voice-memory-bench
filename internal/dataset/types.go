package dataset

import (
	"context"

	"github.com/mahimairaja/vbench/internal/schema"
)

// Loader is the interface every dataset backend implements.
type Loader interface {
	// Name is the short identifier used in configs (e.g. "locomo").
	Name() string
	// Download fetches raw data into cacheDir if not already present.
	Download(ctx context.Context, cacheDir string) error
	// IsCached reports whether raw data is present and verifiable.
	IsCached(cacheDir string) bool
	// Load yields benchmark items. maxItems<=0 means "all".
	Load(cacheDir, subset string, maxItems int) ([]schema.BenchmarkItem, error)
}
