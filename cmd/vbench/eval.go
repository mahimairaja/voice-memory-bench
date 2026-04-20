package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/mahimairaja/vbench/internal/engine"
	"github.com/mahimairaja/vbench/internal/report"
	"github.com/mahimairaja/vbench/internal/schema"
)

func newEvalCmd() *cobra.Command {
	var (
		configPath string
		runID      string
		maxItems   int
		verbose    bool
	)
	cmd := &cobra.Command{
		Use:   "eval",
		Short: "Run a voice-fitness benchmark end-to-end",
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("--config is required")
			}
			cfg, err := loadConfig(configPath)
			if err != nil {
				return err
			}
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()
			res, err := engine.Run(ctx, engine.Options{
				Config:   cfg,
				RunID:    runID,
				MaxItems: maxItems,
				Verbose:  verbose,
			})
			if err != nil {
				return err
			}
			fmt.Println(report.Summary(res.Provider, res.PerLevel))
			fmt.Fprintf(os.Stderr, "\nrun_id=%s  report=%s\n", res.RunID, res.OutputJSON)
			return nil
		},
	}
	cmd.Flags().StringVar(&configPath, "config", "", "Path to run YAML config")
	cmd.Flags().StringVar(&runID, "run-id", "", "Resume run with this ID (default: generate new)")
	cmd.Flags().IntVar(&maxItems, "max-items", 0, "Cap the number of benchmark items (override config)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose progress output")
	return cmd
}

func loadConfig(path string) (*schema.RunConfig, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	var cfg schema.RunConfig
	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	cfg.ApplyEnvExpansion()
	return &cfg, nil
}
