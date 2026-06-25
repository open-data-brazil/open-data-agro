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

func testRegistry(entry catalog.RegistryEntry) *catalog.Registry {
	return catalog.NewRegistry([]catalog.RegistryEntry{entry})
}

func TestPromoteRequiresBronze(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	cfg := config.LakeConfig{
		StorageMode:      config.StorageModeLocal,
		LakeLocalRoot:    root,
		DeltaStoragePath: filepath.Join(root, "silver"),
		DeltaMinVersions: 30,
	}
	reg := testRegistry(catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	})

	promoter := NewPromoter(cfg, reg)
	_, err := promoter.Promote(context.Background(), PromoteOptions{DatasetID: "conab.estimativa-graos"})
	if err == nil {
		t.Fatal("expected error when bronze is missing")
	}
}

func TestPromoteLocalIntegration(t *testing.T) {
	if os.Getenv("DELTA_INTEGRATION") != "1" {
		t.Skip("set DELTA_INTEGRATION=1 with python deltalake installed")
	}

	root := t.TempDir()
	silverRoot := filepath.Join(root, "silver")
	cfg := config.LakeConfig{
		StorageMode:      config.StorageModeLocal,
		LakeLocalRoot:    root,
		DeltaStoragePath: silverRoot,
		DeltaMinVersions: 30,
	}

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}
	reg := testRegistry(entry)

	raw := []byte("col_a;col_b\n1;foo\n2;bar\n")
	parquetBytes, _, err := ingest.ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}

	bronzeDir := filepath.Join(root, "bronze", "conab", "estimativa-graos", "ingest_date=2026-06-25")
	if err := os.MkdirAll(bronzeDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bronzeDir, "part-test.parquet"), parquetBytes, 0o644); err != nil {
		t.Fatalf("WriteFile parquet: %v", err)
	}

	promoter := NewPromoter(cfg, reg)
	result, err := promoter.Promote(context.Background(), PromoteOptions{DatasetID: "conab.estimativa-graos"})
	if err != nil {
		t.Fatalf("Promote: %v", err)
	}
	if result.RowCount != 2 {
		t.Fatalf("row_count: got %d want 2", result.RowCount)
	}
	if _, err := os.Stat(filepath.Join(result.SilverDir, "_delta_log")); err != nil {
		t.Fatalf("delta log missing: %v", err)
	}
}
