// Package main orchestrates bronze → silver Delta promotions and DuckDB jobs.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/db"
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
		Short:   "Promote bronze Parquet to Delta silver and run DuckDB jobs",
		Long:    "Open Data Agro processor — local-first DuckDB reads and Delta silver promotions.",
		Version: version,
	}

	root.AddCommand(newVersionCmd())
	root.AddCommand(newPromoteCmd())
	root.AddCommand(newSmokeCmd())
	root.AddCommand(newQualityCmd())
	root.AddCommand(newSyncPostgresCmd())
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

func newSmokeCmd() *cobra.Command {
	var dataset string
	var ingestDate string

	cmd := &cobra.Command{
		Use:   "smoke",
		Short: "Count bronze Parquet rows via DuckDB",
		Long:  "Runs duckdb/scripts/smoke_read_parquet.sql against local or S3 bronze paths.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if dataset == "" {
				return fmt.Errorf("--dataset is required")
			}
			return runSmoke(cmd, dataset, ingestDate)
		},
	}

	cmd.Flags().StringVar(&dataset, "dataset", "", "Catalog dataset ID (e.g. conab.estimativa-graos)")
	cmd.Flags().StringVar(&ingestDate, "ingest-date", "", "Optional ingest_date partition (YYYY-MM-DD)")
	_ = cmd.MarkFlagRequired("dataset")
	return cmd
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

func newQualityCmd() *cobra.Command {
	var dataset string
	var checkpoint string

	cmd := &cobra.Command{
		Use:   "quality",
		Short: "Run Great Expectations bronze checkpoint",
		Long:  "Validates lake/bronze Parquet against the dataset expectation suite (local files only).",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if dataset == "" {
				return fmt.Errorf("--dataset is required")
			}
			return runQuality(cmd, dataset, checkpoint)
		},
	}

	cmd.Flags().StringVar(&dataset, "dataset", "", "Catalog dataset ID (e.g. conab.estimativa-graos)")
	cmd.Flags().StringVar(&checkpoint, "checkpoint", "", "Checkpoint name (default: mapped from dataset)")
	_ = cmd.MarkFlagRequired("dataset")
	return cmd
}

func runQuality(cmd *cobra.Command, datasetID, checkpoint string) error {
	cfg, err := config.LoadLakeFromEnv()
	if err != nil {
		return err
	}

	reg, err := catalog.LoadDefaultRegistry()
	if err != nil {
		return err
	}

	gate := processor.NewQualityGate(cfg, reg)
	result, err := gate.RunBronzeCheckpoint(cmd.Context(), processor.QualityOptions{
		DatasetID:  datasetID,
		Checkpoint: checkpoint,
	})
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(
		cmd.OutOrStdout(),
		"quality %s checkpoint=%s passed=%d/%d bronze=%s\n",
		result.DatasetID,
		result.Checkpoint,
		result.SuccessfulExpectations,
		result.EvaluatedExpectations,
		result.BronzeDir,
	)
	return err
}

func runSmoke(cmd *cobra.Command, datasetID, ingestDate string) error {
	cfg, err := config.LoadLakeFromEnv()
	if err != nil {
		return err
	}

	reg, err := catalog.LoadDefaultRegistry()
	if err != nil {
		return err
	}

	smoker, err := processor.NewSmoker(cfg, reg)
	if err != nil {
		return err
	}

	result, err := smoker.Smoke(cmd.Context(), processor.SmokeOptions{
		DatasetID:  datasetID,
		IngestDate: ingestDate,
	})
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(
		cmd.OutOrStdout(),
		"smoke %s rows=%d bronze=%s mode=%s\n",
		result.DatasetID,
		result.RowCount,
		result.BronzeURI,
		cfg.StorageMode,
	)
	return err
}

func runPromote(cmd *cobra.Command, datasetID string) error {
	ctx := cmd.Context()
	cfg, err := config.LoadLakeFromEnv()
	if err != nil {
		return err
	}

	reg, err := catalog.LoadDefaultRegistry()
	if err != nil {
		return err
	}

	var repo *db.Repository
	var jobID string
	if cfg.DatabaseURL != "" {
		repo, err = db.NewRepository(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer repo.Close()
		jobID, err = repo.CreatePromotionJob(ctx, datasetID)
		if err != nil {
			if db.ErrPromotionSchemaMissing(err) {
				return fmt.Errorf("promotion audit table missing: run make migrate-up")
			}
			return err
		}
	}

	finish := func(status db.PromotionStatus, result *processor.PromoteResult, promoteErr error) error {
		if repo == nil || jobID == "" {
			return promoteErr
		}
		var rowCount *int
		var silverPath string
		if result != nil {
			rowCount = &result.RowCount
			silverPath = result.SilverDir
		}
		var errMsg *string
		if promoteErr != nil {
			msg := promoteErr.Error()
			errMsg = &msg
		}
		if finishErr := repo.FinishPromotionJob(ctx, jobID, status, rowCount, silverPath, cfg.StorageMode, errMsg); finishErr != nil {
			if promoteErr != nil {
				return fmt.Errorf("finish promotion job: %v (original: %v)", finishErr, promoteErr)
			}
			return finishErr
		}
		return promoteErr
	}

	promoter := processor.NewPromoter(cfg, reg)

	quality := processor.NewQualityGate(cfg, reg)
	if _, err := quality.RunBronzeCheckpoint(ctx, processor.QualityOptions{DatasetID: datasetID}); err != nil {
		if finishErr := finish(db.PromotionQualityFailed, nil, err); finishErr != nil {
			return finishErr
		}
		return err
	}

	result, err := promoter.Promote(ctx, processor.PromoteOptions{DatasetID: datasetID})
	if err != nil {
		status := db.PromotionFailed
		if finishErr := finish(status, nil, err); finishErr != nil {
			return finishErr
		}
		return err
	}

	status := db.PromotionSuccess
	if result.RowCount == 0 {
		status = db.PromotionSkipped
	}
	if err := finish(status, result, nil); err != nil {
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

func newSyncPostgresCmd() *cobra.Command {
	var martFilter string

	cmd := &cobra.Command{
		Use:   "sync-postgres",
		Short: "Sync gold marts into PostgreSQL analytics schema",
		Long:  "Discovers lake/gold/mart_*/mart.parquet files and full-refreshes analytics.* tables in PostgreSQL.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSyncPostgres(cmd, martFilter)
		},
	}

	cmd.Flags().StringVar(&martFilter, "marts", "", "Comma-separated analytics table names to sync (default: all gold marts)")
	return cmd
}

func runSyncPostgres(cmd *cobra.Command, martFilter string) error {
	ctx := cmd.Context()
	cfg, err := config.LoadLakeFromEnv()
	if err != nil {
		return err
	}
	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		return fmt.Errorf("DATABASE_URL is required for sync-postgres")
	}

	repo, err := db.NewRepository(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer repo.Close()

	filter := processor.ParseMartFilter(martFilter)
	if envFilter := strings.TrimSpace(os.Getenv("UNIFIED_DB_SYNC_MARTS")); envFilter != "" && martFilter == "" {
		filter = processor.ParseMartFilter(envFilter)
	}

	syncer, err := processor.NewSyncPostgres(cfg, repo)
	if err != nil {
		return err
	}

	result, err := syncer.Sync(ctx, processor.SyncPostgresOptions{MartFilter: filter})
	if err != nil {
		return err
	}

	for _, table := range result.Tables {
		_, _ = fmt.Fprintf(
			cmd.OutOrStdout(),
			"synced analytics.%s rows=%d gold=%s min=%s max=%s\n",
			table.TableName,
			table.RowCount,
			table.GoldPath,
			table.MinDate,
			table.MaxDate,
		)
	}
	_, err = fmt.Fprintf(
		cmd.OutOrStdout(),
		"sync-postgres run=%s status=%s tables=%d\n",
		result.RunID,
		result.Status,
		len(result.Tables),
	)
	return err
}
