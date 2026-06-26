package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestUSDAPSDSojaGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readUSDAIngestTestdata(t, "psd_soja.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("usda.psd-soja"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 4 {
		t.Fatalf("rowCount: got %d want 4", rowCount)
	}
}

func readUSDAIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "usda", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
