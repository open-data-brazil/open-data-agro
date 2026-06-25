package processor_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/ingest"
	"github.com/open-data-brazil/open-data-agro/internal/processor"
)

func TestDuckDBDeltaScanSmoke(t *testing.T) {
	if os.Getenv("DELTA_INTEGRATION") != "1" {
		t.Skip("set DELTA_INTEGRATION=1 with python deltalake and duckdb installed")
	}
	if _, err := exec.LookPath("duckdb"); err != nil {
		t.Skip("duckdb CLI not installed")
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
	reg := catalog.NewRegistry([]catalog.RegistryEntry{entry})

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

	promoter := processor.NewPromoter(cfg, reg)
	result, err := promoter.Promote(context.Background(), processor.PromoteOptions{DatasetID: "conab.estimativa-graos"})
	if err != nil {
		t.Fatalf("Promote: %v", err)
	}

	query := "INSTALL delta; LOAD delta; SELECT count(*) FROM delta_scan('" + result.SilverDir + "')"
	cmd := exec.CommandContext(context.Background(), "duckdb", "-csv", "-c", query)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("duckdb delta_scan: %v output=%s", err, out)
	}
	if !contains(string(out), "1") {
		t.Fatalf("unexpected count output: %s", out)
	}
}

func contains(haystack, needle string) bool {
	return len(needle) == 0 || (len(haystack) >= len(needle) && (haystack == needle || len(haystack) > 0 && containsSubstr(haystack, needle)))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
