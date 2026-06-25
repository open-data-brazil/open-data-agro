package db

import (
	"context"
	"strings"
	"time"
)

// PromotionStatus is persisted in catalog.promotion_jobs.status.
type PromotionStatus string

const (
	PromotionRunning PromotionStatus = "running"
	PromotionSuccess PromotionStatus = "success"
	PromotionFailed  PromotionStatus = "failed"
	PromotionSkipped PromotionStatus = "skipped"
)

// PromotionRecord is a row from catalog.promotion_jobs.
type PromotionRecord struct {
	ID           string
	DatasetID    string
	StartedAt    time.Time
	FinishedAt   *time.Time
	Status       PromotionStatus
	RowCount     *int
	SilverPath   *string
	StorageMode  *string
	ErrorMessage *string
}

// CreatePromotionJob inserts a running promotion job and returns its id.
func (r *Repository) CreatePromotionJob(ctx context.Context, datasetID string) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `
		INSERT INTO catalog.promotion_jobs (dataset_id, status)
		VALUES ($1, $2)
		RETURNING id
	`, datasetID, PromotionRunning).Scan(&id)
	return id, err
}

// FinishPromotionJob updates terminal promotion status.
func (r *Repository) FinishPromotionJob(
	ctx context.Context,
	jobID string,
	status PromotionStatus,
	rowCount *int,
	silverPath, storageMode string,
	errorMessage *string,
) error {
	var pathPtr *string
	if silverPath != "" {
		pathPtr = &silverPath
	}
	var modePtr *string
	if storageMode != "" {
		modePtr = &storageMode
	}
	_, err := r.pool.Exec(ctx, `
		UPDATE catalog.promotion_jobs
		SET finished_at = NOW(),
		    status = $2,
		    row_count = $3,
		    silver_path = $4,
		    storage_mode = $5,
		    error_message = $6
		WHERE id = $1
	`, jobID, status, rowCount, pathPtr, modePtr, errorMessage)
	return err
}

// ListRecentPromotions returns the latest promotion jobs.
func (r *Repository) ListRecentPromotions(ctx context.Context, limit int) ([]PromotionRecord, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, dataset_id, started_at, finished_at, status, row_count, silver_path, storage_mode, error_message
		FROM catalog.promotion_jobs
		ORDER BY started_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PromotionRecord
	for rows.Next() {
		var rec PromotionRecord
		if scanErr := rows.Scan(
			&rec.ID,
			&rec.DatasetID,
			&rec.StartedAt,
			&rec.FinishedAt,
			&rec.Status,
			&rec.RowCount,
			&rec.SilverPath,
			&rec.StorageMode,
			&rec.ErrorMessage,
		); scanErr != nil {
			return nil, scanErr
		}
		out = append(out, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// ErrPromotionSchemaMissing helps callers detect unmigrated databases.
func ErrPromotionSchemaMissing(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "promotion_jobs") && strings.Contains(err.Error(), "does not exist")
}
