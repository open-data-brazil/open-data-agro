package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// SyncRunStatus is persisted in analytics.sync_runs.status.
type SyncRunStatus string

const (
	SyncRunning SyncRunStatus = "running"
	SyncSuccess SyncRunStatus = "success"
	SyncFailed  SyncRunStatus = "failed"
	SyncPartial SyncRunStatus = "partial"
)

// SyncTableRecord is a row in analytics.sync_tables.
type SyncTableRecord struct {
	TableName    string
	GoldPath     string
	RowCount     int64
	MinDate      string
	MaxDate      string
	ParquetBytes int64
}

// ErrAnalyticsSchemaMissing indicates migration 000005 has not been applied.
func ErrAnalyticsSchemaMissing(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, `schema "analytics"`) ||
		strings.Contains(msg, "analytics.sync_runs") ||
		strings.Contains(msg, "does not exist")
}

// CreateSyncRun inserts analytics.sync_runs and returns the run ID.
func (r *Repository) CreateSyncRun(ctx context.Context, lakeRoot string) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `
		INSERT INTO analytics.sync_runs (status, lake_root)
		VALUES ($1, $2)
		RETURNING id::text
	`, SyncRunning, lakeRoot).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("create sync run: %w", err)
	}
	return id, nil
}

// FinishSyncRun updates analytics.sync_runs with final status.
func (r *Repository) FinishSyncRun(ctx context.Context, runID string, status SyncRunStatus, tablesSynced int, errMsg *string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE analytics.sync_runs
		SET finished_at = NOW(),
		    status = $2,
		    tables_synced = $3,
		    error_message = $4
		WHERE id = $1::uuid
	`, runID, status, tablesSynced, errMsg)
	if err != nil {
		return fmt.Errorf("finish sync run: %w", err)
	}
	return nil
}

// InsertSyncTable records one mart sync in analytics.sync_tables.
func (r *Repository) InsertSyncTable(ctx context.Context, runID string, rec SyncTableRecord) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO analytics.sync_tables (
			run_id, table_name, gold_path, row_count, min_date, max_date, parquet_bytes
		) VALUES ($1::uuid, $2, $3, $4, $5, $6, $7)
	`, runID, rec.TableName, rec.GoldPath, rec.RowCount, nullIfEmpty(rec.MinDate), nullIfEmpty(rec.MaxDate), rec.ParquetBytes)
	if err != nil {
		return fmt.Errorf("insert sync table %s: %w", rec.TableName, err)
	}
	return nil
}

// TableColumnNames returns lowercase column names for analytics.table.
func (r *Repository) TableColumnNames(ctx context.Context, tableName string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = 'analytics' AND table_name = $1
		ORDER BY ordinal_position
	`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		cols = append(cols, name)
	}
	return cols, rows.Err()
}

// ReplaceAnalyticsTable recreates analytics.table with TEXT columns and bulk-loads rows.
func (r *Repository) ReplaceAnalyticsTable(ctx context.Context, tableName string, columns []string, rows [][]string) (int64, error) {
	if len(columns) == 0 {
		return 0, fmt.Errorf("no columns for analytics.%s", tableName)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	quotedTable := pgx.Identifier{"analytics", tableName}.Sanitize()
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", quotedTable)
	if _, err := tx.Exec(ctx, dropSQL); err != nil {
		return 0, fmt.Errorf("drop analytics.%s: %w", tableName, err)
	}

	colDefs := make([]string, len(columns))
	for i, col := range columns {
		colDefs[i] = fmt.Sprintf("%s TEXT", pgx.Identifier{col}.Sanitize())
	}
	createSQL := fmt.Sprintf("CREATE TABLE %s (%s)", quotedTable, strings.Join(colDefs, ", "))
	if _, err := tx.Exec(ctx, createSQL); err != nil {
		return 0, fmt.Errorf("create analytics.%s: %w", tableName, err)
	}

	if len(rows) == 0 {
		if err := tx.Commit(ctx); err != nil {
			return 0, err
		}
		return 0, nil
	}

	pgxCols := make([]string, len(columns))
	for i, col := range columns {
		pgxCols[i] = col
	}

	copyCount, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"analytics", tableName},
		pgxCols,
		stringRowsToAny(rows),
	)
	if err != nil {
		return 0, fmt.Errorf("copy into analytics.%s: %w", tableName, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return copyCount, nil
}

// CreateAnalyticsIndex creates an index when it does not already exist.
func (r *Repository) CreateAnalyticsIndex(ctx context.Context, tableName, indexSuffix string, columns []string) error {
	if len(columns) == 0 {
		return nil
	}
	indexName := fmt.Sprintf("%s%s", tableName, indexSuffix)
	parts := make([]string, len(columns))
	for i, col := range columns {
		parts[i] = pgx.Identifier{col}.Sanitize()
	}
	quotedTable := pgx.Identifier{"analytics", tableName}.Sanitize()
	sql := fmt.Sprintf(
		"CREATE INDEX IF NOT EXISTS %s ON %s (%s)",
		pgx.Identifier{indexName}.Sanitize(),
		quotedTable,
		strings.Join(parts, ", "),
	)
	_, err := r.pool.Exec(ctx, sql)
	return err
}

// CountAnalyticsTable returns row count for analytics.table.
func (r *Repository) CountAnalyticsTable(ctx context.Context, tableName string) (int64, error) {
	quotedTable := pgx.Identifier{"analytics", tableName}.Sanitize()
	var count int64
	err := r.pool.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", quotedTable)).Scan(&count)
	return count, err
}

func nullIfEmpty(v string) any {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}

func stringRowsToAny(rows [][]string) pgx.CopyFromSource {
	return &copyFromStringRows{rows: rows, idx: -1}
}

type copyFromStringRows struct {
	rows [][]string
	idx  int
}

func (c *copyFromStringRows) Next() bool {
	c.idx++
	return c.idx < len(c.rows)
}

func (c *copyFromStringRows) Values() ([]any, error) {
	if c.idx < 0 || c.idx >= len(c.rows) {
		return nil, errors.New("copy row out of range")
	}
	row := c.rows[c.idx]
	out := make([]any, len(row))
	for i, v := range row {
		out[i] = v
	}
	return out, nil
}

func (c *copyFromStringRows) Err() error { return nil }

// SyncStartedAt returns started_at for a sync run (testing helper).
func (r *Repository) SyncStartedAt(ctx context.Context, runID string) (time.Time, error) {
	var started time.Time
	err := r.pool.QueryRow(ctx, `SELECT started_at FROM analytics.sync_runs WHERE id = $1::uuid`, runID).Scan(&started)
	return started, err
}
