package mapa

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSISSERCSVGolden(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("testdata", "sisser_psr_2025.sample.csv"))
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}

	headers, rows, err := parseSISSERCSV(raw, "PSR - 2025")
	if err != nil {
		t.Fatalf("parseSISSERCSV: %v", err)
	}
	if len(headers) < 10 {
		t.Fatalf("headers: got %d", len(headers))
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}
	if rows[0][0] != "PSR - 2025" {
		t.Fatalf("periodo_arquivo: got %q", rows[0][0])
	}
}
