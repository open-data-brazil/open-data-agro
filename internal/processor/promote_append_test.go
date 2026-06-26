package processor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/ingest"
)

func TestPromoteAppendSecondVersion(t *testing.T) {
	if os.Getenv("DELTA_INTEGRATION") != "1" {
		t.Skip("set DELTA_INTEGRATION=1 with python deltalake installed")
	}

	root := t.TempDir()
	cfg := config.LakeConfig{
		StorageMode:      config.StorageModeLocal,
		LakeLocalRoot:    root,
		DeltaStoragePath: filepath.Join(root, "silver"),
		DeltaMinVersions: 30,
	}
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}
	reg := catalog.NewRegistry([]catalog.RegistryEntry{entry})
	promoter := NewPromoter(cfg, reg)

	writeBronze := func(dir string, raw []byte) {
		t.Helper()
		parquetBytes, _, err := ingest.ConvertToParquet(entry, raw)
		if err != nil {
			t.Fatalf("ConvertToParquet: %v", err)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "part-a.parquet"), parquetBytes, 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	writeBronze(
		filepath.Join(root, "bronze", "conab", "estimativa-graos", "ingest_date=2026-06-25"),
		[]byte("a;b\n1;x\n"),
	)
	first, err := promoter.Promote(context.Background(), PromoteOptions{DatasetID: "conab.estimativa-graos"})
	if err != nil {
		t.Fatalf("first promote: %v", err)
	}

	writeBronze(
		filepath.Join(root, "bronze", "conab", "estimativa-graos", "ingest_date=2026-06-26"),
		[]byte("a;b\n2;y\n3;z\n"),
	)
	second, err := promoter.Promote(context.Background(), PromoteOptions{DatasetID: "conab.estimativa-graos"})
	if err != nil {
		t.Fatalf("second promote: %v", err)
	}
	if first.RowCount != 1 || second.RowCount != 2 {
		t.Fatalf("expected 1 then 2 new rows, got %d and %d", first.RowCount, second.RowCount)
	}
	logDir := filepath.Join(first.SilverDir, "_delta_log")
	entries, err := os.ReadDir(logDir)
	if err != nil {
		t.Fatalf("read delta log: %v", err)
	}
	jsonCount := 0
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" && !strings.Contains(e.Name(), ".crc") {
			jsonCount++
		}
	}
	if jsonCount < 2 {
		t.Fatalf("expected at least 2 delta versions, got %d files in %s", jsonCount, logDir)
	}

	duckPath := filepath.Join(t.TempDir(), "delta-version.duckdb")
	duck, err := NewDuckDB(duckPath)
	if err != nil {
		t.Fatalf("NewDuckDB: %v", err)
	}
	ctx := context.Background()
	if _, err := duck.RunSQL(ctx, "FORCE INSTALL delta; LOAD delta;"); err != nil {
		t.Fatalf("load delta extension: %v", err)
	}
	silverPath := first.SilverDir
	v0Out, err := duck.RunSQL(ctx, "SELECT count(*) AS row_count FROM delta_scan('"+silverPath+"', version => 0)")
	if err != nil {
		t.Fatalf("delta_scan version 0: %v", err)
	}
	v0Count, err := parseCountCSV(v0Out)
	if err != nil {
		t.Fatalf("parse version 0 count: %v (out=%q)", err, v0Out)
	}
	if v0Count != 1 {
		t.Fatalf("version 0 row_count: got %d want 1", v0Count)
	}

	latestOut, err := duck.RunSQL(ctx, "SELECT count(*) AS row_count FROM delta_scan('"+silverPath+"')")
	if err != nil {
		t.Fatalf("delta_scan latest: %v", err)
	}
	latestCount, err := parseCountCSV(latestOut)
	if err != nil {
		t.Fatalf("parse latest count: %v (out=%q)", err, latestOut)
	}
	if latestCount != 3 {
		t.Fatalf("latest row_count: got %d want 3", latestCount)
	}
}
