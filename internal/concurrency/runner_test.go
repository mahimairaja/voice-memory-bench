package concurrency

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync/atomic"
	"testing"
	"time"
)

func TestRun_preservesOrder(t *testing.T) {
	ctx := context.Background()
	got := Run(ctx, 5, 3, func(_ context.Context, i int) (int, error) {
		return i * 10, nil
	})
	if len(got) != 5 {
		t.Fatalf("want 5 results, got %d", len(got))
	}
	for i, r := range got {
		if r.Index != i {
			t.Errorf("result[%d].Index = %d", i, r.Index)
		}
		if r.Value != i*10 {
			t.Errorf("result[%d].Value = %d", i, r.Value)
		}
	}
}

func TestRun_respectsConcurrencyCap(t *testing.T) {
	ctx := context.Background()
	var inFlight, peak atomic.Int64
	Run(ctx, 20, 3, func(_ context.Context, _ int) (int, error) {
		cur := inFlight.Add(1)
		defer inFlight.Add(-1)
		for {
			p := peak.Load()
			if cur <= p || peak.CompareAndSwap(p, cur) {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
		return 0, nil
	})
	if peak.Load() > 3 {
		t.Errorf("peak in-flight = %d, want <= 3", peak.Load())
	}
}

func TestRun_errorPropagates(t *testing.T) {
	ctx := context.Background()
	got := Run(ctx, 3, 2, func(_ context.Context, i int) (int, error) {
		if i == 1 {
			return 0, errors.New("boom")
		}
		return i, nil
	})
	if got[1].Err == nil || got[0].Err != nil || got[2].Err != nil {
		t.Errorf("unexpected error pattern: %+v", got)
	}
}

func TestRun_zeroTasks(t *testing.T) {
	ctx := context.Background()
	if got := Run(ctx, 0, 4, func(_ context.Context, i int) (int, error) {
		return i, nil
	}); len(got) != 0 {
		t.Errorf("want empty result for n=0, got %d", len(got))
	}
}

func TestPercentile_empty(t *testing.T) {
	if got := Percentile(nil, 95); got != 0 {
		t.Errorf("Percentile(nil, 95) = %v", got)
	}
}

func TestPercentile_single(t *testing.T) {
	if got := Percentile([]float64{42}, 50); got != 42 {
		t.Errorf("Percentile([42],50) = %v", got)
	}
}

func TestPercentile_known(t *testing.T) {
	xs := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	cases := []struct {
		p    float64
		want float64
	}{
		{0, 1},
		{50, 5.5},
		{95, 9.55},
		{99, 9.91},
		{100, 10},
	}
	for _, c := range cases {
		got := Percentile(xs, c.p)
		if math.Abs(got-c.want) > 1e-6 {
			t.Errorf("Percentile(p=%v) = %v, want %v", c.p, got, c.want)
		}
	}
}

func TestPercentile_unsorted(t *testing.T) {
	xs := []float64{5, 3, 9, 1, 7}
	got := Percentile(xs, 50)
	if got != 5 {
		t.Errorf("median of unsorted = %v, want 5", got)
	}
}

func TestRun_scaleSmokeConsistent(t *testing.T) {
	ctx := context.Background()
	got := Run(ctx, 100, 8, func(_ context.Context, i int) (string, error) {
		return fmt.Sprintf("v%d", i), nil
	})
	for i, r := range got {
		want := fmt.Sprintf("v%d", i)
		if r.Value != want {
			t.Fatalf("result[%d] = %q, want %q", i, r.Value, want)
		}
	}
}
