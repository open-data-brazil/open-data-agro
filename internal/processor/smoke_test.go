package processor

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/ingest"
)

func TestSmokeLocalIntegration(t *testing.T) {
	if os.Getenv("DUCKDB_INTEGRATION") != "1" {
		t.Skip("set DUCKDB_INTEGRATION=1 with duckdb CLI installed")
	}

	root := t.TempDir()
	cfg := config.LakeConfig{
		StorageMode:   config.StorageModeLocal,
		LakeLocalRoot: root,
		DuckDBPath:    ":memory:",
	}
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}
	reg := catalog.NewRegistry([]catalog.RegistryEntry{entry})

	parquetBytes, _, err := ingest.ConvertToParquet(entry, []byte("a;b\n1;x\n2;y\n"))
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	bronzeDir := filepath.Join(root, "bronze", "conab", "estimativa-graos", "ingest_date=2026-06-25")
	if err := os.MkdirAll(bronzeDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bronzeDir, "part-a.parquet"), parquetBytes, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	smoker, err := NewSmoker(cfg, reg)
	if err != nil {
		t.Fatalf("NewSmoker: %v", err)
	}
	result, err := smoker.Smoke(context.Background(), SmokeOptions{DatasetID: "conab.estimativa-graos"})
	if err != nil {
		t.Fatalf("Smoke: %v", err)
	}
	if result.RowCount != 2 {
		t.Fatalf("row_count: got %d want 2", result.RowCount)
	}
}

func TestPreviewPromoteLocalIntegration(t *testing.T) {
	if os.Getenv("DUCKDB_INTEGRATION") != "1" {
		t.Skip("set DUCKDB_INTEGRATION=1 with duckdb CLI installed")
	}

	root := t.TempDir()
	cfg := config.LakeConfig{
		StorageMode:      config.StorageModeLocal,
		LakeLocalRoot:    root,
		DeltaStoragePath: filepath.Join(root, "silver"),
		DuckDBPath:         ":memory:",
	}
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}
	parquetBytes, _, err := ingest.ConvertToParquet(entry, []byte("a;b\n1;x\n"))
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	bronzeDir := filepath.Join(root, "bronze", "conab", "estimativa-graos", "ingest_date=2026-06-25")
	if err := os.MkdirAll(bronzeDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bronzeDir, "part-a.parquet"), parquetBytes, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	count, err := PreviewPromote(context.Background(), cfg, "conab.estimativa-graos", "")
	if err != nil {
		t.Fatalf("PreviewPromote: %v", err)
	}
	if count != 1 {
		t.Fatalf("preview count: got %d want 1", count)
	}
}
