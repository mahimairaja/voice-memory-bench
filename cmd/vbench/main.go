package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is set at build time with -ldflags "-X main.version=..."
var version = "0.1.0-dev"

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vbench",
		Short: "vbench — voice-agent fitness benchmark for self-hostable memory frameworks",
		Long: `vbench measures whether a memory framework is shippable inside a voice agent.

The headline output is one line per concurrency level:

    mem0 @ 4x: ACCEPTABLE (p95 = 380 ms, quality = 0.72, tokens = 420)

Three verdicts map to the voice latency budget:

    p95 < 300 ms           → EXCELLENT
    300 <= p95 <= 500 ms   → ACCEPTABLE
    p95 > 500 ms           → FAIL
`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	cmd.AddCommand(newEvalCmd())
	cmd.AddCommand(newDatasetsCmd())
	cmd.AddCommand(newProvidersCmd())
	return cmd
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
