package main

import (
	"context"
	"fmt"
	"os"

	"github.com/open-data-brazil/open-data-agro/internal/alerts"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/db"
	"github.com/open-data-brazil/open-data-agro/internal/ingest"
	"github.com/open-data-brazil/open-data-agro/internal/storage"
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
		Use:     "ingestor",
		Short:   "Download official agro datasets and land them in the data lake",
		Long:    "Open Data Agro ingestor — CONAB and other official Brazilian agricultural data sources.",
		Version: version,
	}

	root.AddCommand(newVersionCmd())
	root.AddCommand(newCatalogCmd())
	root.AddCommand(newRunCmd())
	root.AddCommand(newStatusCmd())
	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print ingestor version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "ingestor %s\n", version)
			return err
		},
	}
}

func newCatalogCmd() *cobra.Command {
	catalogCmd := &cobra.Command{
		Use:   "catalog",
		Short: "Inspect the dataset catalog registry",
	}

	catalogCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List registered datasets",
		RunE: func(cmd *cobra.Command, _ []string) error {
			reg, err := catalog.LoadDefaultRegistry()
			if err != nil {
				return err
			}
			for _, entry := range reg.Entries() {
				line := fmt.Sprintf("%s\t%s\t%s\t%s", entry.DatasetID, entry.ConabSection, entry.Schedule, entry.Format)
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), line); err != nil {
					return err
				}
			}
			return nil
		},
	})

	return catalogCmd
}

func newRunCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "run <dataset_id>",
		Short: "Download, convert, and land a dataset to bronze",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			cfg, err := config.LoadFromEnv()
			if err != nil {
				return err
			}

			reg, err := catalog.LoadDefaultRegistry()
			if err != nil {
				return err
			}

			repo, err := db.NewRepository(ctx, cfg.DatabaseURL)
			if err != nil {
				return err
			}
			defer repo.Close()

			store, err := storage.NewBronzeStore(cfg)
			if err != nil {
				return err
			}

			runner := ingest.NewRunner(reg, repo, store, alerts.New(cfg.AlertWebhookURL))
			result, err := runner.Run(ctx, ingest.RunOptions{
				DatasetID: args[0],
				DryRun:    dryRun,
			})
			if err != nil {
				return err
			}

			msg := fmt.Sprintf("status=%s sha256=%s backend=%s", result.Status, result.SHA256, result.Backend)
			if result.R2Key != "" {
				msg += fmt.Sprintf(" key=%s", result.R2Key)
			}
			if result.Skipped {
				msg += " (skipped: identical checksum already ingested)"
			}
			_, printErr := fmt.Fprintln(cmd.OutOrStdout(), msg)
			return printErr
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Resolve and download without landing parquet")
	return cmd
}

func newStatusCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show recent ingest jobs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			cfg, err := config.LoadFromEnv()
			if err != nil {
				return err
			}

			repo, err := db.NewRepository(ctx, cfg.DatabaseURL)
			if err != nil {
				return err
			}
			defer repo.Close()

			jobs, err := repo.ListRecentJobs(ctx, limit)
			if err != nil {
				return err
			}

			for _, job := range jobs {
				finished := "-"
				if job.FinishedAt != nil {
					finished = job.FinishedAt.UTC().Format("2006-01-02T15:04:05Z")
				}
				errMsg := ""
				if job.ErrorMessage != nil {
					errMsg = *job.ErrorMessage
				}
				line := fmt.Sprintf(
					"%s\t%s\t%s\t%s\t%t\t%s",
					job.StartedAt.UTC().Format("2006-01-02T15:04:05Z"),
					finished,
					job.DatasetID,
					job.Status,
					job.DryRun,
					errMsg,
				)
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), line); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&limit, "last", 10, "Number of recent jobs to show")
	return cmd
}
