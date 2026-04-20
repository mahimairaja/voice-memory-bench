package concurrency

import (
	"context"
	"math"
	"sort"
	"sync"
)

// Task is one unit of work the runner will invoke. It receives a zero-based
// index so the implementation can correlate results.
type Task[R any] func(ctx context.Context, index int) (R, error)

// Result is the per-task outcome preserved in original input order.
type Result[R any] struct {
	Index int
	Value R
	Err   error
}

// Run executes n tasks against fn with at most `concurrency` in flight at once.
// Returned results are sorted by input index so callers can match input/output
// positionally. The first error does not cancel remaining work — the engine
// decides whether to treat per-task failures as fatal.
func Run[R any](ctx context.Context, n int, concurrency int, fn Task[R]) []Result[R] {
	if concurrency < 1 {
		concurrency = 1
	}
	if n == 0 {
		return nil
	}
	results := make([]Result[R], n)
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		i := i
		// Acquire a slot, but bail out if ctx is cancelled while we're waiting
		// for in-flight workers to release one. Without this, a cancelled run
		// can block forever on a full semaphore.
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			results[i] = Result[R]{Index: i, Err: ctx.Err()}
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			if ctx.Err() != nil {
				results[i] = Result[R]{Index: i, Err: ctx.Err()}
				return
			}
			v, err := fn(ctx, i)
			results[i] = Result[R]{Index: i, Value: v, Err: err}
		}()
	}
	wg.Wait()
	sort.Slice(results, func(a, b int) bool { return results[a].Index < results[b].Index })
	return results
}

// Percentile returns the p-th percentile of xs (p in [0,100]), using linear
// interpolation between the closest ranks. Returns 0 if xs is empty.
func Percentile(xs []float64, p float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	if len(xs) == 1 {
		return xs[0]
	}
	sorted := append([]float64(nil), xs...)
	sort.Float64s(sorted)
	if p <= 0 {
		return sorted[0]
	}
	if p >= 100 {
		return sorted[len(sorted)-1]
	}
	rank := (p / 100.0) * float64(len(sorted)-1)
	lo := int(math.Floor(rank))
	hi := int(math.Ceil(rank))
	if lo == hi {
		return sorted[lo]
	}
	frac := rank - float64(lo)
	return sorted[lo] + frac*(sorted[hi]-sorted[lo])
}
