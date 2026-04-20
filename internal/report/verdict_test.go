package report

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mahimairaja/vbench/internal/schema"
)

func TestClassifyP95_boundaries(t *testing.T) {
	cases := []struct {
		p95  float64
		want Verdict
	}{
		{0, VerdictExcellent},
		{299.999, VerdictExcellent},
		{300, VerdictAcceptable},
		{500, VerdictAcceptable},
		{500.001, VerdictFail},
		{5000, VerdictFail},
	}
	for _, c := range cases {
		got := ClassifyP95(c.p95)
		if got != c.want {
			t.Errorf("ClassifyP95(%v) = %q, want %q", c.p95, got, c.want)
		}
	}
}

func TestHeadline_format(t *testing.T) {
	ms := schema.MemScore{
		Quality:           0.72,
		LatencyP95Ms:      380,
		TokenFootprintP50: 420,
		Concurrency:       4,
	}
	got := Headline("mem0", ms)
	want := "mem0 @ 4x: ACCEPTABLE (p95 = 380 ms, quality = 0.72, tokens = 420)"
	if got != want {
		t.Errorf("\n got: %q\nwant: %q", got, want)
	}
}

func TestSummary_oneLinePerLevel(t *testing.T) {
	levels := []schema.MemScore{
		{Quality: 0.74, LatencyP95Ms: 210, TokenFootprintP50: 390, Concurrency: 1},
		{Quality: 0.72, LatencyP95Ms: 380, TokenFootprintP50: 420, Concurrency: 4},
	}
	got := Summary("mem0", levels)
	parts := strings.Split(got, "\n")
	if len(parts) != 2 {
		t.Fatalf("want 2 lines, got %d: %q", len(parts), got)
	}
	if !strings.Contains(parts[0], "EXCELLENT") {
		t.Errorf("line 0 should be EXCELLENT: %q", parts[0])
	}
	if !strings.Contains(parts[1], "ACCEPTABLE") {
		t.Errorf("line 1 should be ACCEPTABLE: %q", parts[1])
	}
}

func TestWriteJSON_roundtrips(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "memscore.json")
	levels := []schema.MemScore{{
		Quality:           0.5,
		LatencyP50Ms:      150,
		LatencyP95Ms:      280,
		LatencyP99Ms:      400,
		TokenFootprintP50: 300,
		Concurrency:       1,
		NumQuestions:      10,
	}}
	if err := WriteJSON(path, "mem0", "locomo", "run-x", levels); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var rep Report
	if err := json.Unmarshal(data, &rep); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if rep.Provider != "mem0" || rep.Dataset != "locomo" || rep.RunID != "run-x" {
		t.Errorf("top-level fields wrong: %+v", rep)
	}
	if len(rep.Levels) != 1 || rep.Levels[0].Verdict != VerdictExcellent {
		t.Errorf("levels wrong: %+v", rep.Levels)
	}
	if len(rep.Headline) != 1 || !strings.Contains(rep.Headline[0], "EXCELLENT") {
		t.Errorf("headline wrong: %+v", rep.Headline)
	}
}
