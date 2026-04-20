package schema

import (
	"fmt"
	"os"
	"regexp"
)

// RunConfig is the top-level YAML config for one benchmark run.
type RunConfig struct {
	RunName     string         `yaml:"run_name"`
	Dataset     DatasetConfig  `yaml:"dataset"`
	Provider    ProviderConfig `yaml:"provider"`
	AnswerLLM   LLMConfig      `yaml:"answer_llm"`
	JudgeLLM    LLMConfig      `yaml:"judge_llm"`
	Concurrency []int          `yaml:"concurrency"`
	Seed        int            `yaml:"seed"`
	OutputDir   string         `yaml:"output_dir"`
}

// DatasetConfig selects a benchmark dataset.
type DatasetConfig struct {
	Name     string `yaml:"name"`
	Subset   string `yaml:"subset,omitempty"`
	Path     string `yaml:"path,omitempty"`
	MaxItems int    `yaml:"max_items,omitempty"`
}

// ProviderConfig selects a memory provider and its sidecar config.
type ProviderConfig struct {
	Name    string                 `yaml:"name"`
	Command []string               `yaml:"command,omitempty"`
	Config  map[string]interface{} `yaml:"config"`
}

// LLMConfig configures a single LLM endpoint.
type LLMConfig struct {
	Model       string  `yaml:"model"`
	BaseURL     string  `yaml:"base_url,omitempty"`
	APIKeyEnv   string  `yaml:"api_key_env,omitempty"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
	Seed        int     `yaml:"seed"`
}

// Validate enforces MVP invariants.
func (c *RunConfig) Validate() error {
	if c.RunName == "" {
		return fmt.Errorf("run_name is required")
	}
	if c.Dataset.Name == "" {
		return fmt.Errorf("dataset.name is required")
	}
	if c.Dataset.Name == "custom" && c.Dataset.Path == "" {
		return fmt.Errorf("dataset.path is required when dataset.name=custom")
	}
	if c.Provider.Name == "" {
		return fmt.Errorf("provider.name is required")
	}
	if c.AnswerLLM.Model == "" {
		return fmt.Errorf("answer_llm.model is required")
	}
	if c.JudgeLLM.Model == "" {
		return fmt.Errorf("judge_llm.model is required")
	}
	if len(c.Concurrency) == 0 {
		c.Concurrency = []int{1}
	}
	for _, n := range c.Concurrency {
		if n < 1 {
			return fmt.Errorf("concurrency levels must be >= 1, got %d", n)
		}
	}
	if c.OutputDir == "" {
		c.OutputDir = "runs"
	}
	return nil
}

var envPattern = regexp.MustCompile(`\$\{([A-Z0-9_]+)\}`)

// ExpandEnv recursively expands ${VAR} references in string-valued config entries.
func ExpandEnv(v interface{}) interface{} {
	switch t := v.(type) {
	case string:
		return envPattern.ReplaceAllStringFunc(t, func(match string) string {
			name := match[2 : len(match)-1]
			return os.Getenv(name)
		})
	case map[string]interface{}:
		for k, val := range t {
			t[k] = ExpandEnv(val)
		}
		return t
	case []interface{}:
		for i, val := range t {
			t[i] = ExpandEnv(val)
		}
		return t
	default:
		return v
	}
}

// ApplyEnvExpansion expands env vars across the config where it matters.
func (c *RunConfig) ApplyEnvExpansion() {
	if c.Provider.Config != nil {
		ExpandEnv(c.Provider.Config)
	}
}
