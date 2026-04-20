package schema

import (
	"os"
	"testing"
)

func TestValidate_requiresRunName(t *testing.T) {
	c := &RunConfig{}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing run_name")
	}
}

func TestValidate_fillsDefaults(t *testing.T) {
	c := &RunConfig{
		RunName:   "x",
		Dataset:   DatasetConfig{Name: "locomo"},
		Provider:  ProviderConfig{Name: "mem0"},
		AnswerLLM: LLMConfig{Model: "gpt-4o-mini"},
		JudgeLLM:  LLMConfig{Model: "gpt-4o-mini"},
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.OutputDir; got != "runs" {
		t.Errorf("OutputDir default = %q, want %q", got, "runs")
	}
	if len(c.Concurrency) != 1 || c.Concurrency[0] != 1 {
		t.Errorf("Concurrency default = %v, want [1]", c.Concurrency)
	}
}

func TestValidate_rejectsZeroConcurrency(t *testing.T) {
	c := &RunConfig{
		RunName:     "x",
		Dataset:     DatasetConfig{Name: "locomo"},
		Provider:    ProviderConfig{Name: "mem0"},
		AnswerLLM:   LLMConfig{Model: "m"},
		JudgeLLM:    LLMConfig{Model: "m"},
		Concurrency: []int{0},
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for concurrency=0")
	}
}

func TestValidate_customDatasetNeedsPath(t *testing.T) {
	c := &RunConfig{
		RunName:   "x",
		Dataset:   DatasetConfig{Name: "custom"},
		Provider:  ProviderConfig{Name: "mem0"},
		AnswerLLM: LLMConfig{Model: "m"},
		JudgeLLM:  LLMConfig{Model: "m"},
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for custom dataset without path")
	}
}

func TestExpandEnv_recurses(t *testing.T) {
	t.Setenv("VBENCH_TEST_SECRET", "hunter2")
	cfg := map[string]interface{}{
		"postgres_url": "postgres://u:${VBENCH_TEST_SECRET}@h/d",
		"nested": map[string]interface{}{
			"key": "value-${VBENCH_TEST_SECRET}",
		},
	}
	ExpandEnv(cfg)
	got := cfg["postgres_url"].(string)
	if got != "postgres://u:hunter2@h/d" {
		t.Errorf("postgres_url = %q", got)
	}
	nested := cfg["nested"].(map[string]interface{})
	if nested["key"].(string) != "value-hunter2" {
		t.Errorf("nested.key = %q", nested["key"])
	}
}

func TestApplyEnvExpansion_noop(t *testing.T) {
	_ = os.Setenv("NOT_PRESENT", "")
	c := &RunConfig{Provider: ProviderConfig{Config: nil}}
	c.ApplyEnvExpansion()
}
