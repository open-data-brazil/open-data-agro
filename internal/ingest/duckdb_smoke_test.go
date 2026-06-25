package ingest_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/ingest"
)

func TestDuckDBReadParquetSmoke(t *testing.T) {
	if _, err := exec.LookPath("duckdb"); err != nil {
		t.Skip("duckdb CLI not installed")
	}

	root := t.TempDir()
	parquetPath := filepath.Join(root, "sample.parquet")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}
	raw := []byte("a;b\n1;x\n2;y\n")
	parquetBytes, rowCount, err := ingest.ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 2 {
		t.Fatalf("rowCount: got %d want 2", rowCount)
	}
	if err := os.WriteFile(parquetPath, parquetBytes, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "duckdb", "-csv", "-c", "SELECT count(*) FROM read_parquet('"+parquetPath+"')")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("duckdb: %v output=%s", err, out)
	}
	if !containsLine(out, "2") {
		t.Fatalf("unexpected duckdb output: %s", out)
	}
}

func containsLine(data []byte, want string) bool {
	for _, line := range splitLines(data) {
		if line == want {
			return true
		}
	}
	return false
}

func splitLines(data []byte) []string {
	var lines []string
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, string(data[start:i]))
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, string(data[start:]))
	}
	return lines
}
