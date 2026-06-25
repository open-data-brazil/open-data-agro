package processor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/lake"
)

// SmokeOptions controls a DuckDB bronze read smoke test.
type SmokeOptions struct {
	DatasetID  string
	IngestDate string // optional YYYY-MM-DD partition filter
}

// SmokeResult summarizes a smoke read.
type SmokeResult struct {
	DatasetID string
	RowCount  int
	BronzeURI string
}

// Smoker runs DuckDB validation queries against bronze Parquet.
type Smoker struct {
	cfg      config.LakeConfig
	registry *catalog.Registry
	duckdb   *DuckDB
}

// NewSmoker wires smoke dependencies.
func NewSmoker(cfg config.LakeConfig, registry *catalog.Registry) (*Smoker, error) {
	duck, err := NewDuckDB(cfg.DuckDBPath)
	if err != nil {
		return nil, err
	}
	return &Smoker{cfg: cfg, registry: registry, duckdb: duck}, nil
}

// Smoke counts rows in bronze Parquet for a dataset.
func (s *Smoker) Smoke(ctx context.Context, opts SmokeOptions) (*SmokeResult, error) {
	if _, err := s.registry.Require(opts.DatasetID); err != nil {
		return nil, err
	}

	bronzeURI, err := s.bronzeURI(opts)
	if err != nil {
		return nil, err
	}

	scriptPath, err := ScriptPath("smoke_read_parquet.sql")
	if err != nil {
		return nil, err
	}

	vars := map[string]string{"bronze_uri": bronzeURI}
	sql, err := loadScript(scriptPath, vars)
	if err != nil {
		return nil, err
	}

	out, err := s.duckdb.RunSQL(ctx, s.extensionSetup()+sql)
	if err != nil {
		return nil, fmt.Errorf("smoke read parquet: %w", err)
	}

	rowCount, err := parseCountCSV(out)
	if err != nil {
		return nil, fmt.Errorf("parse smoke output %q: %w", out, err)
	}

	return &SmokeResult{
		DatasetID: opts.DatasetID,
		RowCount:  rowCount,
		BronzeURI: bronzeURI,
	}, nil
}

func (s *Smoker) bronzeURI(opts SmokeOptions) (string, error) {
	lakeRoot := lake.NormalizeRoot(s.cfg.LakeLocalRoot)

	switch s.cfg.StorageMode {
	case config.StorageModeMinIO, config.StorageModeR2:
		slug, err := lake.DatasetSlug(opts.DatasetID)
		if err != nil {
			return "", err
		}
		bucket := s.cfg.MinIOBucket
		if s.cfg.StorageMode == config.StorageModeR2 {
			bucket = s.cfg.R2Bucket
		}
		prefix := fmt.Sprintf("s3://%s/bronze/conab/%s", bucket, slug)
		if opts.IngestDate != "" {
			return prefix + "/ingest_date=" + opts.IngestDate + "/*.parquet", nil
		}
		return prefix + "/**/*.parquet", nil
	default:
		slug, err := lake.DatasetSlug(opts.DatasetID)
		if err != nil {
			return "", err
		}
		base := filepath.Join(lakeRoot, "bronze", "conab", slug)
		if opts.IngestDate != "" {
			return filepath.Join(base, "ingest_date="+opts.IngestDate, "*.parquet"), nil
		}
		return filepath.Join(base, "**", "*.parquet"), nil
	}
}

func (s *Smoker) extensionSetup() string {
	switch s.cfg.StorageMode {
	case config.StorageModeMinIO:
		return S3ExtensionsSQL(map[string]string{
			"access_key": s.cfg.MinIOAccessKey,
			"secret_key": s.cfg.MinIOSecretKey,
			"endpoint":   strings.TrimPrefix(s.cfg.MinIOEndpoint, "http://"),
			"url_style":  "path",
			"use_ssl":    "false",
		})
	case config.StorageModeR2:
		endpoint := strings.TrimPrefix(strings.TrimPrefix(s.cfg.R2Endpoint, "https://"), "http://")
		return S3ExtensionsSQL(map[string]string{
			"access_key": s.cfg.R2AccessKeyID,
			"secret_key": s.cfg.R2SecretAccessKey,
			"endpoint":   endpoint,
			"url_style":  "vhost",
			"use_ssl":    "true",
		})
	default:
		return LocalExtensionsSQL()
	}
}

func loadScript(path string, vars map[string]string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read script %s: %w", path, err)
	}
	return SubstituteVars(string(data), vars), nil
}

func parseCountCSV(out string) (int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("expected header and count line")
	}
	return strconv.Atoi(strings.TrimSpace(lines[1]))
}
