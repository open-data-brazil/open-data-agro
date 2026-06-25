package db

import (
	"context"
	"time"
)

// LatestIngestRecord is a row from catalog.v_latest_successful_ingest.
type LatestIngestRecord struct {
	DatasetID  string
	JobID      string
	StartedAt  time.Time
	FinishedAt *time.Time
	Status     JobStatus
	DryRun     bool
}

// ListLatestSuccessfulIngest returns the most recent successful job per dataset.
func (r *Repository) ListLatestSuccessfulIngest(ctx context.Context) ([]LatestIngestRecord, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT dataset_id, job_id, started_at, finished_at, status, dry_run
		FROM catalog.v_latest_successful_ingest
		ORDER BY dataset_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []LatestIngestRecord
	for rows.Next() {
		var rec LatestIngestRecord
		if scanErr := rows.Scan(
			&rec.DatasetID,
			&rec.JobID,
			&rec.StartedAt,
			&rec.FinishedAt,
			&rec.Status,
			&rec.DryRun,
		); scanErr != nil {
			return nil, scanErr
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

// ListFailedJobsLast7d returns failed ingest jobs from the last seven days.
func (r *Repository) ListFailedJobsLast7d(ctx context.Context) ([]JobRecord, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, dataset_id, started_at, finished_at, status, error_message, dry_run
		FROM catalog.v_failed_jobs_last_7d
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []JobRecord
	for rows.Next() {
		var rec JobRecord
		if scanErr := rows.Scan(
			&rec.ID,
			&rec.DatasetID,
			&rec.StartedAt,
			&rec.FinishedAt,
			&rec.Status,
			&rec.ErrorMessage,
			&rec.DryRun,
		); scanErr != nil {
			return nil, scanErr
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}
