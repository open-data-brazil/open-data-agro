package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// JobStatus is persisted in catalog.ingest_jobs.status.
type JobStatus string

const (
	JobRunning JobStatus = "running"
	JobSuccess JobStatus = "success"
	JobFailed  JobStatus = "failed"
	JobSkipped JobStatus = "skipped"
)

// JobRecord is a row from catalog.ingest_jobs.
type JobRecord struct {
	ID           string
	DatasetID    string
	StartedAt    time.Time
	FinishedAt   *time.Time
	Status       JobStatus
	ErrorMessage *string
	DryRun       bool
}

// FileRecord is a row from catalog.ingest_files.
type FileRecord struct {
	ID            string
	JobID         string
	DatasetID     string
	SHA256        string
	RowCount      *int
	R2Key         *string
	ContentType   *string
	LastModified  *string
	FileSizeBytes *int64
}

// Repository persists ingest audit metadata in PostgreSQL.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository connects to PostgreSQL using DATABASE_URL.
func NewRepository(ctx context.Context, databaseURL string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &Repository{pool: pool}, nil
}

// Close releases the connection pool.
func (r *Repository) Close() {
	r.pool.Close()
}

// UpsertDatasetRegistry syncs a catalog entry into catalog.dataset_registry.
func (r *Repository) UpsertDatasetRegistry(ctx context.Context, entry catalog.RegistryEntry) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO catalog.dataset_registry (
			dataset_id, source_url, source_portal_url, format, schedule, conab_section, portal_label, discovered_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (dataset_id) DO UPDATE SET
			source_url = EXCLUDED.source_url,
			source_portal_url = EXCLUDED.source_portal_url,
			format = EXCLUDED.format,
			schedule = EXCLUDED.schedule,
			conab_section = EXCLUDED.conab_section,
			portal_label = EXCLUDED.portal_label,
			discovered_at = EXCLUDED.discovered_at,
			updated_at = NOW()
	`,
		entry.DatasetID.String(),
		entry.SourceURL,
		entry.PortalURL(),
		string(entry.Format),
		entry.Schedule,
		entry.ConabSection,
		entry.PortalLabel,
		entry.DiscoveredAt,
	)
	return err
}

// FindFileBySHA256 returns an existing ingest file for idempotency checks.
func (r *Repository) FindFileBySHA256(ctx context.Context, datasetID, sha256 string) (*FileRecord, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, job_id, dataset_id, sha256, row_count, r2_key, content_type, last_modified, file_size_bytes
		FROM catalog.ingest_files
		WHERE dataset_id = $1 AND sha256 = $2
	`, datasetID, sha256)

	var rec FileRecord
	err := row.Scan(
		&rec.ID,
		&rec.JobID,
		&rec.DatasetID,
		&rec.SHA256,
		&rec.RowCount,
		&rec.R2Key,
		&rec.ContentType,
		&rec.LastModified,
		&rec.FileSizeBytes,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

// CreateJob inserts a running ingest job and returns its id.
func (r *Repository) CreateJob(ctx context.Context, datasetID string, dryRun bool) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `
		INSERT INTO catalog.ingest_jobs (dataset_id, status, dry_run)
		VALUES ($1, $2, $3)
		RETURNING id
	`, datasetID, JobRunning, dryRun).Scan(&id)
	return id, err
}

// FinishJob updates terminal job status.
func (r *Repository) FinishJob(ctx context.Context, jobID string, status JobStatus, errorMessage *string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE catalog.ingest_jobs
		SET finished_at = NOW(), status = $2, error_message = $3
		WHERE id = $1
	`, jobID, status, errorMessage)
	return err
}

// CreateIngestFile records a landed parquet object.
func (r *Repository) CreateIngestFile(ctx context.Context, rec FileRecord) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO catalog.ingest_files (
			job_id, dataset_id, sha256, row_count, r2_key, content_type, last_modified, file_size_bytes, discovered_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`,
		rec.JobID,
		rec.DatasetID,
		rec.SHA256,
		rec.RowCount,
		rec.R2Key,
		rec.ContentType,
		rec.LastModified,
		rec.FileSizeBytes,
	)
	return err
}

// DeleteIngestFile removes metadata for rollback after a failed upload.
func (r *Repository) DeleteIngestFile(ctx context.Context, datasetID, sha256 string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM catalog.ingest_files
		WHERE dataset_id = $1 AND sha256 = $2
	`, datasetID, sha256)
	return err
}

// ListRecentJobs returns the latest ingest jobs.
func (r *Repository) ListRecentJobs(ctx context.Context, limit int) ([]JobRecord, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, dataset_id, started_at, finished_at, status, error_message, dry_run
		FROM catalog.ingest_jobs
		ORDER BY started_at DESC
		LIMIT $1
	`, limit)
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

func isNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
