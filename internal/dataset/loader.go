package dataset

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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

// Names returns the registered dataset names in deterministic (sorted) order.
func Names() []string {
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	sort.Strings(out)
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

var unsafeIDChars = regexp.MustCompile(`[^A-Za-z0-9._-]`)

// SafeID normalises a dataset-supplied identifier so it is safe to use as a
// filesystem path component. Loaders must run untrusted identifiers through
// SafeID before storing them on BenchmarkItem.ItemID or EvaluationQuestion
// identifiers — the engine uses those values as filenames verbatim.
func SafeID(id string) string {
	s := strings.TrimSpace(id)
	s = unsafeIDChars.ReplaceAllString(s, "_")
	if s == "" || s == "." || s == ".." || strings.HasPrefix(s, "..") {
		return "_"
	}
	return s
}
