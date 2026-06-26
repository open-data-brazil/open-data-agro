package processor

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/db"
	"github.com/open-data-brazil/open-data-agro/internal/lake"
)

// GoldMart describes one gold mart parquet target.
type GoldMart struct {
	DirName     string
	TableName   string
	ParquetPath string
}

// SyncPostgresOptions controls a gold → PostgreSQL sync run.
type SyncPostgresOptions struct {
	MartFilter []string
}

// TableSyncResult summarizes one mart sync.
type TableSyncResult struct {
	TableName    string
	GoldPath     string
	RowCount     int64
	MinDate      string
	MaxDate      string
	ParquetBytes int64
}

// SyncPostgresResult summarizes a full sync run.
type SyncPostgresResult struct {
	RunID  string
	Tables []TableSyncResult
	Status db.SyncRunStatus
}

// SyncPostgres mirrors gold marts into PostgreSQL analytics schema.
type SyncPostgres struct {
	cfg    config.LakeConfig
	duckdb *DuckDB
	repo   *db.Repository
}

// NewSyncPostgres wires gold → PostgreSQL sync dependencies.
func NewSyncPostgres(cfg config.LakeConfig, repo *db.Repository) (*SyncPostgres, error) {
	duck, err := NewDuckDB(cfg.DuckDBPath)
	if err != nil {
		return nil, err
	}
	return &SyncPostgres{cfg: cfg, duckdb: duck, repo: repo}, nil
}

// Sync discovers gold marts and loads them into PostgreSQL.
func (s *SyncPostgres) Sync(ctx context.Context, opts SyncPostgresOptions) (*SyncPostgresResult, error) {
	lakeRoot := lake.NormalizeRoot(s.cfg.LakeLocalRoot)
	marts, err := DiscoverGoldMarts(lakeRoot, opts.MartFilter)
	if err != nil {
		return nil, err
	}
	if len(marts) == 0 {
		return nil, fmt.Errorf("no gold marts found under %s/gold", lakeRoot)
	}

	runID, err := s.repo.CreateSyncRun(ctx, lakeRoot)
	if err != nil {
		if db.ErrAnalyticsSchemaMissing(err) {
			return nil, fmt.Errorf("analytics schema missing: run make migrate-up")
		}
		return nil, err
	}

	result := &SyncPostgresResult{RunID: runID, Status: db.SyncSuccess}
	var syncErr error
	var failed int

	for _, mart := range marts {
		tableResult, err := s.syncMart(ctx, mart)
		if err != nil {
			failed++
			syncErr = err
			break
		}
		if err := s.repo.InsertSyncTable(ctx, runID, db.SyncTableRecord{
			TableName:    tableResult.TableName,
			GoldPath:     tableResult.GoldPath,
			RowCount:     tableResult.RowCount,
			MinDate:      tableResult.MinDate,
			MaxDate:      tableResult.MaxDate,
			ParquetBytes: tableResult.ParquetBytes,
		}); err != nil {
			failed++
			syncErr = err
			break
		}
		result.Tables = append(result.Tables, *tableResult)
	}

	var errMsg *string
	if syncErr != nil {
		msg := syncErr.Error()
		errMsg = &msg
		if len(result.Tables) > 0 {
			result.Status = db.SyncPartial
		} else {
			result.Status = db.SyncFailed
		}
	}
	if failed > 0 && len(result.Tables) == 0 {
		result.Status = db.SyncFailed
	}

	if finishErr := s.repo.FinishSyncRun(ctx, runID, result.Status, len(result.Tables), errMsg); finishErr != nil {
		if syncErr != nil {
			return result, fmt.Errorf("finish sync run: %v (original: %v)", finishErr, syncErr)
		}
		return result, finishErr
	}
	if syncErr != nil {
		return result, syncErr
	}
	return result, nil
}

func (s *SyncPostgres) syncMart(ctx context.Context, mart GoldMart) (*TableSyncResult, error) {
	info, err := os.Stat(mart.ParquetPath)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", mart.ParquetPath, err)
	}

	columns, rows, err := s.readMartCSV(ctx, mart.ParquetPath)
	if err != nil {
		return nil, fmt.Errorf("export %s: %w", mart.TableName, err)
	}

	rowCount, err := s.repo.ReplaceAnalyticsTable(ctx, mart.TableName, columns, rows)
	if err != nil {
		return nil, err
	}

	if err := s.ensureJoinIndexes(ctx, mart.TableName, columns); err != nil {
		return nil, err
	}

	minDate, maxDate := dateRangeFromRows(columns, rows)
	return &TableSyncResult{
		TableName:    mart.TableName,
		GoldPath:     mart.ParquetPath,
		RowCount:     rowCount,
		MinDate:      minDate,
		MaxDate:      maxDate,
		ParquetBytes: info.Size(),
	}, nil
}

func (s *SyncPostgres) readMartCSV(ctx context.Context, parquetPath string) ([]string, [][]string, error) {
	tempCSV, err := os.CreateTemp("", "sync-postgres-*.csv")
	if err != nil {
		return nil, nil, err
	}
	tempPath := tempCSV.Name()
	_ = tempCSV.Close()
	defer func() { _ = os.Remove(tempPath) }()

	escapedParquet := escapeSQLString(parquetPath)
	escapedCSV := escapeSQLString(tempPath)
	exportSQL := fmt.Sprintf(
		`COPY (SELECT * FROM read_parquet('%s')) TO '%s' (HEADER, DELIMITER ',', QUOTE '"')`,
		escapedParquet,
		escapedCSV,
	)
	if _, err := s.duckdb.RunSQL(ctx, exportSQL); err != nil {
		return nil, nil, err
	}

	file, err := os.Open(tempPath)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	reader.ReuseRecord = true
	header, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("read csv header: %w", err)
	}
	columns := normalizeColumns(header)

	var rows [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("read csv row: %w", err)
		}
		if len(record) < len(columns) {
			padded := make([]string, len(columns))
			copy(padded, record)
			rows = append(rows, padded)
			continue
		}
		if len(record) > len(columns) {
			record = record[:len(columns)]
		}
		row := make([]string, len(columns))
		copy(row, record)
		rows = append(rows, row)
	}
	return columns, rows, nil
}

func (s *SyncPostgres) ensureJoinIndexes(ctx context.Context, tableName string, columns []string) error {
	colSet := make(map[string]struct{}, len(columns))
	for _, col := range columns {
		colSet[strings.ToLower(col)] = struct{}{}
	}
	has := func(name string) bool {
		_, ok := colSet[strings.ToLower(name)]
		return ok
	}

	type indexSpec struct {
		suffix  string
		columns []string
	}
	specs := []indexSpec{
		{suffix: "_cod_ibge_idx", columns: []string{"cod_ibge"}},
		{suffix: "_codigo_ibge_idx", columns: []string{"codigo_ibge"}},
		{suffix: "_produto_safra_idx", columns: []string{"produto", "safra"}},
		{suffix: "_refmonth_idx", columns: []string{"refmonth"}},
		{suffix: "_data_preco_idx", columns: []string{"data_preco"}},
		{suffix: "_capturado_em_idx", columns: []string{"capturado_em"}},
	}

	for _, spec := range specs {
		ok := true
		for _, col := range spec.columns {
			if !has(col) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		if err := s.repo.CreateAnalyticsIndex(ctx, tableName, spec.suffix, spec.columns); err != nil {
			return fmt.Errorf("create index on analytics.%s: %w", tableName, err)
		}
	}
	return nil
}

// DiscoverGoldMarts lists mart_*/mart.parquet under lake/gold.
func DiscoverGoldMarts(lakeRoot string, filter []string) ([]GoldMart, error) {
	goldRoot := filepath.Join(lakeRoot, "gold")
	entries, err := os.ReadDir(goldRoot)
	if err != nil {
		return nil, fmt.Errorf("read gold dir: %w", err)
	}

	filterSet := make(map[string]struct{}, len(filter))
	for _, name := range filter {
		if trimmed := strings.TrimSpace(name); trimmed != "" {
			filterSet[trimmed] = struct{}{}
		}
	}

	var marts []GoldMart
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "mart_") {
			continue
		}
		tableName, err := MartTableName(entry.Name())
		if err != nil {
			continue
		}
		if len(filterSet) > 0 {
			if _, ok := filterSet[tableName]; !ok {
				continue
			}
		}
		parquetPath := filepath.Join(goldRoot, entry.Name(), "mart.parquet")
		if _, err := os.Stat(parquetPath); err != nil {
			continue
		}
		marts = append(marts, GoldMart{
			DirName:     entry.Name(),
			TableName:   tableName,
			ParquetPath: parquetPath,
		})
	}

	sort.Slice(marts, func(i, j int) bool {
		return marts[i].TableName < marts[j].TableName
	})
	return marts, nil
}

// MartTableName maps mart directory name to analytics table name.
func MartTableName(dirName string) (string, error) {
	name := strings.TrimSpace(dirName)
	if !strings.HasPrefix(name, "mart_") {
		return "", fmt.Errorf("invalid mart dir %q", dirName)
	}
	name = strings.TrimPrefix(name, "mart_")
	name = strings.ReplaceAll(name, "__", "_")
	if name == "" {
		return "", fmt.Errorf("empty mart table for %q", dirName)
	}
	return name, nil
}

func normalizeColumns(header []string) []string {
	out := make([]string, len(header))
	seen := make(map[string]int, len(header))
	for i, raw := range header {
		col := strings.TrimSpace(raw)
		if col == "" {
			col = fmt.Sprintf("col_%d", i+1)
		}
		key := strings.ToLower(col)
		if n, ok := seen[key]; ok {
			seen[key] = n + 1
			col = fmt.Sprintf("%s_%d", col, n+1)
		} else {
			seen[key] = 1
		}
		out[i] = col
	}
	return out
}

func dateRangeFromRows(columns []string, rows [][]string) (minDate, maxDate string) {
	dateCols := []string{"refmonth", "data_preco", "capturado_em", "ano_mes", "data", "refyear"}
	colIndex := make(map[string]int, len(columns))
	for i, col := range columns {
		colIndex[strings.ToLower(col)] = i
	}

	var values []string
	for _, name := range dateCols {
		idx, ok := colIndex[name]
		if !ok {
			continue
		}
		for _, row := range rows {
			if idx >= len(row) {
				continue
			}
			v := strings.TrimSpace(row[idx])
			if v != "" {
				values = append(values, v)
			}
		}
	}
	if len(values) == 0 {
		return "", ""
	}
	sort.Strings(values)
	return values[0], values[len(values)-1]
}

func escapeSQLString(raw string) string {
	return strings.ReplaceAll(raw, "'", "''")
}

// ParseMartFilter splits comma-separated mart table names.
func ParseMartFilter(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
