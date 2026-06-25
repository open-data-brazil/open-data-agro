// Package main orchestrates bronze → silver Delta promotions and DuckDB jobs.
package main

import (
	"fmt"
	"os"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/processor"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "processor",
		Short:   "Promote bronze Parquet to Delta silver and run lake jobs",
		Long:    "Open Data Agro processor — local-first Delta Lake promotions via delta-rs (Python).",
		Version: version,
	}

	root.AddCommand(newVersionCmd())
	root.AddCommand(newPromoteCmd())
	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print processor version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "processor %s\n", version)
			return err
		},
	}
}

func newPromoteCmd() *cobra.Command {
	var dataset string

	cmd := &cobra.Command{
		Use:   "promote",
		Short: "Promote bronze Parquet to Delta silver",
		Long:  "Reads lake/bronze Parquet partitions and appends a new Delta version under lake/silver/.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if dataset == "" {
				return fmt.Errorf("--dataset is required")
			}
			return runPromote(cmd, dataset)
		},
	}

	cmd.Flags().StringVar(&dataset, "dataset", "", "Catalog dataset ID (e.g. conab.estimativa-graos)")
	_ = cmd.MarkFlagRequired("dataset")
	return cmd
}

func runPromote(cmd *cobra.Command, datasetID string) error {
	cfg, err := config.LoadLakeFromEnv()
	if err != nil {
		return err
	}

	reg, err := catalog.LoadDefaultRegistry()
	if err != nil {
		return err
	}

	promoter := processor.NewPromoter(cfg, reg)
	result, err := promoter.Promote(cmd.Context(), processor.PromoteOptions{DatasetID: datasetID})
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(
		cmd.OutOrStdout(),
		"promoted %s rows=%d silver=%s mode=%s\n",
		result.DatasetID,
		result.RowCount,
		result.SilverDir,
		result.StorageMode,
	)
	return err
}
