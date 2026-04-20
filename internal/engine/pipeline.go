package engine

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mahimairaja/vbench/internal/dataset"
	"github.com/mahimairaja/vbench/internal/llm"
	"github.com/mahimairaja/vbench/internal/report"
	"github.com/mahimairaja/vbench/internal/schema"
	"github.com/mahimairaja/vbench/internal/sidecar"
)

// Pipeline executes the ingest → index → search → answer → evaluate stages.
// Each stage writes a `.complete` sentinel so subsequent runs with the same
// run_id can resume instead of redoing finished work.
type Pipeline struct {
	Config      *schema.RunConfig
	RunID       string
	ArtifactDir string
	Items       []schema.BenchmarkItem
	Side        *sidecar.Process
	AnswerLLM   *llm.Client
	JudgeLLM    *llm.Client
	MaxItems    int
	Verbose     bool
}

// PipelineResult bundles the per-concurrency MemScores that the report layer
// renders as voice verdicts.
type PipelineResult struct {
	RunID      string
	Provider   string
	Dataset    string
	PerLevel   []schema.MemScore
	OutputJSON string
}

// Options controls pipeline construction.
type Options struct {
	Config   *schema.RunConfig
	RunID    string
	MaxItems int
	Verbose  bool
}

// NewRunID returns a short hex ID, suitable as a directory name.
func NewRunID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return time.Now().Format("20060102-150405") + "-" + hex.EncodeToString(b)
}

// Run executes the whole pipeline end-to-end. It assumes the sidecar has been
// started by the caller and that the run config has been validated.
func Run(ctx context.Context, opts Options) (*PipelineResult, error) {
	cfg := opts.Config
	runID := opts.RunID
	if runID == "" {
		runID = NewRunID()
	}
	artifactDir := filepath.Join(cfg.OutputDir, runID)
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir artifact dir: %w", err)
	}
	if err := writeManifest(artifactDir, cfg, runID); err != nil {
		return nil, err
	}

	// Dataset
	ds, err := dataset.Get(cfg.Dataset.Name)
	if err != nil {
		return nil, err
	}
	cacheDir := dataset.DefaultCacheDir()
	if !ds.IsCached(cacheDir) {
		return nil, fmt.Errorf("dataset %s not cached; run `vbench datasets download %s`", cfg.Dataset.Name, cfg.Dataset.Name)
	}
	max := opts.MaxItems
	if max <= 0 {
		max = cfg.Dataset.MaxItems
	}
	items, err := ds.Load(cacheDir, cfg.Dataset.Subset, max)
	if err != nil {
		return nil, fmt.Errorf("load dataset: %w", err)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("dataset %s yielded zero items", cfg.Dataset.Name)
	}

	// Sidecar
	side, err := startSidecar(ctx, cfg)
	if err != nil {
		return nil, err
	}
	defer side.Shutdown()

	caps, err := side.Client().Capabilities(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch sidecar capabilities: %w", err)
	}
	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "sidecar ready: %s v%s modes=%v\n", caps.ProviderName, caps.ProviderVersion, caps.SupportedRetrievalModes)
	}

	// Answer + judge LLM clients. The engine is a boundary: surface a clear
	// error if either configured api_key_env is unset so the user isn't left
	// to decode a downstream 401 from the model provider. A BaseURL pointing
	// at a local/self-hosted endpoint without auth is still allowed via an
	// empty api_key_env.
	ansKey, err := requireAPIKey("answer_llm", cfg.AnswerLLM)
	if err != nil {
		return nil, err
	}
	jdgKey, err := requireAPIKey("judge_llm", cfg.JudgeLLM)
	if err != nil {
		return nil, err
	}
	answerLLM := llm.New(cfg.AnswerLLM, ansKey)
	judgeLLM := llm.New(cfg.JudgeLLM, jdgKey)

	pipe := &Pipeline{
		Config:      cfg,
		RunID:       runID,
		ArtifactDir: artifactDir,
		Items:       items,
		Side:        side,
		AnswerLLM:   answerLLM,
		JudgeLLM:    judgeLLM,
		MaxItems:    max,
		Verbose:     opts.Verbose,
	}

	if err := pipe.ingest(ctx); err != nil {
		return nil, fmt.Errorf("ingest: %w", err)
	}
	if err := pipe.index(ctx); err != nil {
		return nil, fmt.Errorf("index: %w", err)
	}

	// Search stage runs once per concurrency level.
	var perLevel []schema.MemScore
	for _, c := range cfg.Concurrency {
		if err := pipe.search(ctx, c); err != nil {
			return nil, fmt.Errorf("search @%dx: %w", c, err)
		}
		if err := pipe.answer(ctx, c); err != nil {
			return nil, fmt.Errorf("answer @%dx: %w", c, err)
		}
		ms, err := pipe.evaluate(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("evaluate @%dx: %w", c, err)
		}
		perLevel = append(perLevel, ms)
	}

	// Write final JSON report.
	out := filepath.Join(artifactDir, "memscore.json")
	if err := report.WriteJSON(out, cfg.Provider.Name, cfg.Dataset.Name, runID, perLevel); err != nil {
		return nil, err
	}

	return &PipelineResult{
		RunID:      runID,
		Provider:   cfg.Provider.Name,
		Dataset:    cfg.Dataset.Name,
		PerLevel:   perLevel,
		OutputJSON: out,
	}, nil
}

func startSidecar(ctx context.Context, cfg *schema.RunConfig) (*sidecar.Process, error) {
	cmd := cfg.Provider.Command
	workingDir := ""
	if len(cmd) == 0 {
		switch cfg.Provider.Name {
		case "mem0":
			cmd = sidecar.Mem0DefaultCommand
			workingDir = sidecar.Mem0DefaultWorkingDir
		default:
			return nil, fmt.Errorf("no default sidecar command for provider %q; set provider.command in config", cfg.Provider.Name)
		}
	}
	return sidecar.Spawn(ctx, sidecar.SpawnOptions{
		Command:        cmd,
		WorkingDir:     workingDir,
		ProviderConfig: cfg.Provider.Config,
		ReadyTimeout:   60 * time.Second,
	})
}

func writeManifest(artifactDir string, cfg *schema.RunConfig, runID string) error {
	manifest := map[string]interface{}{
		"run_id":    runID,
		"provider":  cfg.Provider.Name,
		"dataset":   cfg.Dataset.Name,
		"run_name":  cfg.RunName,
		"seed":      cfg.Seed,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	return writeJSON(filepath.Join(artifactDir, "manifest.json"), manifest)
}

func writeJSON(path string, v interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func stageDir(artifactDir, stage string) string {
	return filepath.Join(artifactDir, stage)
}

func completeSentinel(dir string) string { return filepath.Join(dir, ".complete") }

func stageComplete(dir string) bool {
	_, err := os.Stat(completeSentinel(dir))
	return err == nil
}

func markComplete(dir string) error {
	f, err := os.Create(completeSentinel(dir))
	if err != nil {
		return err
	}
	return f.Close()
}

func requireAPIKey(role string, cfg schema.LLMConfig) (string, error) {
	if cfg.APIKeyEnv == "" {
		// No env var name set — caller has declared this endpoint does not
		// require auth (typical for a self-hosted vLLM/Ollama behind BaseURL).
		return "", nil
	}
	v := os.Getenv(cfg.APIKeyEnv)
	if v == "" {
		return "", fmt.Errorf("%s.api_key_env=%q is not set in the environment", role, cfg.APIKeyEnv)
	}
	return v, nil
}
