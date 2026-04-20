package main

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/mahimairaja/vbench/internal/dataset"
)

func newDatasetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datasets",
		Short: "Manage benchmark datasets",
	}
	cmd.AddCommand(newDatasetsListCmd())
	cmd.AddCommand(newDatasetsDownloadCmd())
	cmd.AddCommand(newDatasetsInfoCmd())
	return cmd
}

func newDatasetsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available datasets",
		RunE: func(cmd *cobra.Command, args []string) error {
			names := dataset.Names()
			sort.Strings(names)
			for _, n := range names {
				fmt.Println(n)
			}
			return nil
		},
	}
}

func newDatasetsDownloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "download <name>",
		Short: "Download a dataset into the user cache",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := dataset.Get(args[0])
			if err != nil {
				return err
			}
			cache := dataset.DefaultCacheDir()
			if l.IsCached(cache) {
				fmt.Printf("%s already cached at %s\n", args[0], cache)
				return nil
			}
			if err := l.Download(cache); err != nil {
				return err
			}
			fmt.Printf("downloaded %s → %s\n", args[0], cache)
			return nil
		},
	}
}

func newDatasetsInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <name>",
		Short: "Show dataset metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := dataset.Get(args[0])
			if err != nil {
				return err
			}
			cache := dataset.DefaultCacheDir()
			fmt.Printf("name:   %s\n", l.Name())
			fmt.Printf("cache:  %s\n", cache)
			fmt.Printf("cached: %v\n", l.IsCached(cache))
			return nil
		},
	}
}
