package main

import (
	"fmt"
	"os"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "ingestor",
		Short: "Download official agro datasets and land them in the data lake",
		Long:  "Open Data Agro ingestor — CONAB and other official Brazilian agricultural data sources.",
	}

	root.AddCommand(newCatalogCmd())
	return root
}

func newCatalogCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "catalog",
		Short: "Inspect the dataset catalog registry",
		RunE: func(cmd *cobra.Command, _ []string) error {
			reg, err := catalog.LoadDefaultRegistry()
			if err != nil {
				return err
			}
			for _, id := range reg.ListIDs() {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), id); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
