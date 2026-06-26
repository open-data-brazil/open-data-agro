package processor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DuckDB runs SQL via the DuckDB CLI (no CGO).
type DuckDB struct {
	Binary string
	DBPath string
}

// NewDuckDB resolves the CLI binary and database path from config/env.
func NewDuckDB(dbPath string) (*DuckDB, error) {
	binary := strings.TrimSpace(os.Getenv("DUCKDB_BIN"))
	if binary == "" {
		var err error
		binary, err = exec.LookPath("duckdb")
		if err != nil {
			return nil, fmt.Errorf("duckdb CLI not found (set DUCKDB_BIN or run make duckdb-install): %w", err)
		}
	}
	if strings.TrimSpace(dbPath) == "" {
		dbPath = ":memory:"
	}
	return &DuckDB{Binary: binary, DBPath: dbPath}, nil
}

// RunSQL executes SQL and returns combined stdout.
func (d *DuckDB) RunSQL(ctx context.Context, sql string) (string, error) {
	cmd := exec.CommandContext(ctx, d.Binary, d.DBPath, "-csv", "-c", sql)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("duckdb: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(string(out)), nil
}

// RunScriptFile loads a SQL file, substitutes ${var} placeholders, and executes it.
func (d *DuckDB) RunScriptFile(ctx context.Context, scriptPath string, vars map[string]string) (string, error) {
	data, err := os.ReadFile(scriptPath)
	if err != nil {
		return "", fmt.Errorf("read script %s: %w", scriptPath, err)
	}
	sql := SubstituteVars(string(data), vars)
	return d.RunSQL(ctx, sql)
}

// SubstituteVars replaces ${key} placeholders in SQL templates.
func SubstituteVars(sql string, vars map[string]string) string {
	out := sql
	for key, value := range vars {
		out = strings.ReplaceAll(out, "${"+key+"}", value)
	}
	return out
}

// ScriptPath resolves a path under duckdb/scripts/.
func ScriptPath(name string) (string, error) {
	if v := strings.TrimSpace(os.Getenv("DUCKDB_SCRIPTS_DIR")); v != "" {
		return filepath.Join(v, name), nil
	}
	root, err := findModuleRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "duckdb", "scripts", name), nil
}

// LocalExtensionsSQL returns setup SQL for local lake paths (parquet is built-in).
func LocalExtensionsSQL() string {
	return ""
}

// S3ExtensionsSQL installs httpfs for S3-compatible bronze reads (MinIO/R2).
func S3ExtensionsSQL(vars map[string]string) string {
	return SubstituteVars(`
FORCE INSTALL httpfs;
LOAD httpfs;
CREATE OR REPLACE SECRET lake_s3 (
  TYPE s3,
  KEY_ID '${access_key}',
  SECRET '${secret_key}',
  ENDPOINT '${endpoint}',
  URL_STYLE '${url_style}',
  USE_SSL ${use_ssl}
);
`, vars)
}
