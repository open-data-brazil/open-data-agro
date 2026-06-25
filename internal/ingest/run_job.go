package ingest

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/alerts"
	"github.com/open-data-brazil/open-data-agro/internal/anp"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/conab"
	"github.com/open-data-brazil/open-data-agro/internal/db"
	"github.com/open-data-brazil/open-data-agro/internal/storage"
)

// RunOptions controls a single dataset ingest execution.
type RunOptions struct {
	DatasetID string
	DryRun    bool
}

// RunResult summarizes an ingest attempt.
type RunResult struct {
	JobID    string
	Status   db.JobStatus
	SHA256   string
	R2Key    string
	RowCount int
	Skipped  bool
	Backend  string
}

// Runner orchestrates download, convert, upload, and audit.
type Runner struct {
	registry *catalog.Registry
	repo     *db.Repository
	store    storage.BronzeStore
	conab    *conab.Client
	anp      *anp.Client
	alerts   *alerts.Notifier
}

// NewRunner wires ingest dependencies.
func NewRunner(registry *catalog.Registry, repo *db.Repository, store storage.BronzeStore, notifier *alerts.Notifier) *Runner {
	return &Runner{
		registry: registry,
		repo:     repo,
		store:    store,
		conab:    conab.NewClient(),
		anp:      anp.NewClient(),
		alerts:   notifier,
	}
}

// Run executes one ingest job for a dataset.
func (r *Runner) Run(ctx context.Context, opts RunOptions) (*RunResult, error) {
	entry, err := r.registry.Require(opts.DatasetID)
	if err != nil {
		return nil, err
	}

	if err := r.repo.UpsertDatasetRegistry(ctx, entry); err != nil {
		return nil, fmt.Errorf("sync dataset registry: %w", err)
	}

	jobID, err := r.repo.CreateJob(ctx, opts.DatasetID, opts.DryRun)
	if err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}

	finish := func(status db.JobStatus, message *string, result *RunResult, runErr error) (*RunResult, error) {
		if finishErr := r.repo.FinishJob(ctx, jobID, status, message); finishErr != nil {
			return result, fmt.Errorf("finish job: %w (original: %v)", finishErr, runErr)
		}
		if status == db.JobFailed {
			r.alerts.Notify(ctx, slog.LevelError, fmt.Sprintf("ingest failed for %s", opts.DatasetID), "job_id", jobID, "error", derefString(message))
		}
		return result, runErr
	}

	sourceURL, err := ResolveSourceURL(entry)
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	download, err := DownloadSource(ctx, entry, r.conab, r.anp)
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	ingestDate := time.Now().UTC()
	fp := NewFingerprint(download.Body, download.ContentType, download.LastModified, 0, ingestDate)

	existing, err := r.repo.FindFileBySHA256(ctx, opts.DatasetID, fp.SHA256)
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}
	if existing != nil {
		result := &RunResult{
			JobID:   jobID,
			Status:  db.JobSkipped,
			SHA256:  fp.SHA256,
			Skipped: true,
			Backend: r.store.Backend(),
		}
		if existing.R2Key != nil {
			result.R2Key = *existing.R2Key
		}
		if existing.RowCount != nil {
			result.RowCount = *existing.RowCount
		}
		_, _ = finish(db.JobSkipped, nil, result, nil)
		return result, nil
	}

	if opts.DryRun {
		result := &RunResult{
			JobID:   jobID,
			Status:  db.JobSuccess,
			SHA256:  fp.SHA256,
			Backend: r.store.Backend(),
		}
		_, _ = finish(db.JobSuccess, nil, result, nil)
		return result, nil
	}

	staged, err := StageRaw(jobID, download.Body)
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}
	defer staged.Cleanup()

	parquetBytes, rowCount, err := ConvertToParquetFromFile(entry, staged.Path)
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}
	fp.RowCount = rowCount

	bronzeKey, err := BronzeKey(opts.DatasetID, ingestDate, fp.PartID)
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	if err := r.store.Put(ctx, bronzeKey, parquetBytes, "application/vnd.apache.parquet"); err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	metadataKey, err := PartitionMetadataKey(opts.DatasetID, ingestDate)
	if err != nil {
		_ = r.store.Delete(ctx, bronzeKey)
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}
	meta := NewBronzeMetadata(entry, fp, sourceURL, r.store.Backend())
	metaBytes, err := MarshalBronzeMetadata(meta)
	if err != nil {
		_ = r.store.Delete(ctx, bronzeKey)
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}
	if err := r.store.Put(ctx, metadataKey, metaBytes, "application/json"); err != nil {
		_ = r.store.Delete(ctx, bronzeKey)
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	rowCountCopy := rowCount
	fileSize := fp.FileSizeBytes
	fileRec := db.FileRecord{
		JobID:         jobID,
		DatasetID:     opts.DatasetID,
		SHA256:        fp.SHA256,
		RowCount:      &rowCountCopy,
		R2Key:         &bronzeKey,
		ContentType:   &download.ContentType,
		LastModified:  &download.LastModified,
		FileSizeBytes: &fileSize,
	}
	if err := r.repo.CreateIngestFile(ctx, fileRec); err != nil {
		if delErr := r.store.Delete(ctx, bronzeKey); delErr != nil {
			msg := fmt.Sprintf("metadata insert failed and orphan cleanup failed: %v / %v", err, delErr)
			_, _ = finish(db.JobFailed, &msg, nil, err)
			return nil, err
		}
		_ = r.store.Delete(ctx, metadataKey)
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	result := &RunResult{
		JobID:    jobID,
		Status:   db.JobSuccess,
		SHA256:   fp.SHA256,
		R2Key:    bronzeKey,
		RowCount: rowCount,
		Backend:  r.store.Backend(),
	}
	_, _ = finish(db.JobSuccess, nil, result, nil)
	return result, nil
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
