package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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
	var crop string
	var fromYear int
	var toYear int
	var year int
	var ufList string

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
			var ufs []string
			if strings.TrimSpace(ufList) != "" {
				for _, part := range strings.Split(ufList, ",") {
					code := strings.TrimSpace(part)
					if code != "" {
						ufs = append(ufs, code)
					}
				}
			}

			result, err := runner.Run(ctx, ingest.RunOptions{
				DatasetID: args[0],
				DryRun:    dryRun,
				Crop:      crop,
				FromYear:  fromYear,
				ToYear:    toYear,
				Year:      year,
				UFs:       ufs,
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
	cmd.Flags().StringVar(&crop, "crop", "", "PAM crop filter (soja, milho, trigo, or all)")
	cmd.Flags().IntVar(&fromYear, "from", 0, "PAM start year (inclusive)")
	cmd.Flags().IntVar(&toYear, "to", 0, "PAM end year (inclusive)")
	cmd.Flags().IntVar(&year, "year", 0, "INMET BDMEP year for annual ZIP pulls")
	cmd.Flags().StringVar(&ufList, "uf", "", "Filter by UF (comma-separated; numeric codes for IBGE PAM, state abbreviations for INMET)")
	return cmd
}

func newStatusCmd() *cobra.Command {
	var limit int
	var failed bool
	var latest bool

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

			if failed {
				jobs, listErr := repo.ListFailedJobsLast7d(ctx)
				if listErr != nil {
					return listErr
				}
				return printJobs(cmd, jobs)
			}

			if latest {
				records, listErr := repo.ListLatestSuccessfulIngest(ctx)
				if listErr != nil {
					return listErr
				}
				for _, rec := range records {
					finished := "-"
					if rec.FinishedAt != nil {
						finished = rec.FinishedAt.UTC().Format("2006-01-02T15:04:05Z")
					}
					line := fmt.Sprintf(
						"%s\t%s\t%s\t%s\t%t\t%s",
						rec.StartedAt.UTC().Format("2006-01-02T15:04:05Z"),
						finished,
						rec.DatasetID,
						rec.Status,
						rec.DryRun,
						rec.JobID,
					)
					if _, err := fmt.Fprintln(cmd.OutOrStdout(), line); err != nil {
						return err
					}
				}
				return nil
			}

			jobs, err := repo.ListRecentJobs(ctx, limit)
			if err != nil {
				return err
			}
			return printJobs(cmd, jobs)
		},
	}

	cmd.Flags().IntVar(&limit, "last", 10, "Number of recent jobs to show")
	cmd.Flags().BoolVar(&failed, "failed", false, "Show failed jobs from the last 7 days (uses v_failed_jobs_last_7d)")
	cmd.Flags().BoolVar(&latest, "latest", false, "Show latest successful ingest per dataset (uses v_latest_successful_ingest)")
	return cmd
}

func printJobs(cmd *cobra.Command, jobs []db.JobRecord) error {
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
}
