package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestB3FuturoSojaGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readB3IngestTestdata(t, "futuro_soja.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("b3.futuro-soja"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func readB3IngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "b3", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
