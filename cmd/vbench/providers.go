package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// providerEntry is what `vbench providers list` prints. The MVP ships Mem0
// only; new entries are added alongside their sidecar package.
type providerEntry struct {
	Name        string
	SidecarDir  string
	Backing     string
	Status      string
	Description string
}

var mvpProviders = []providerEntry{
	{
		Name:        "mem0",
		SidecarDir:  "sidecars/mem0",
		Backing:     "Postgres + pgvector",
		Status:      "supported",
		Description: "Mem0 OSS — semantic/hybrid retrieval, pgvector store.",
	},
}

func newProvidersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Short: "Manage memory providers",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List MVP-supported providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%-10s  %-22s  %-10s  %s\n", "NAME", "BACKING", "STATUS", "DESCRIPTION")
			for _, p := range mvpProviders {
				fmt.Printf("%-10s  %-22s  %-10s  %s\n", p.Name, p.Backing, p.Status, p.Description)
			}
			return nil
		},
	})
	return cmd
}
