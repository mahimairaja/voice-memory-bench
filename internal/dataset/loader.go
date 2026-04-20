package dataset

import (
	"fmt"
	"os"
	"path/filepath"
)

// Registry maps dataset names to their loader implementations.
var registry = map[string]Loader{
	"locomo": &LoCoMo{},
}

// Get returns the loader for the given dataset name.
func Get(name string) (Loader, error) {
	l, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown dataset %q (available: %v)", name, Names())
	}
	return l, nil
}

// Names returns the list of registered dataset names.
func Names() []string {
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	return out
}

// DefaultCacheDir returns the user-level cache directory for datasets.
// Falls back to ./datasets/cache if the user cache dir is unavailable.
func DefaultCacheDir() string {
	if dir, err := os.UserCacheDir(); err == nil {
		return filepath.Join(dir, "vbench", "datasets")
	}
	return filepath.Join("datasets", "cache")
}
