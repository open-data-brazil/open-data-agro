package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestANTTPracasPedagioGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readANTTTestdata(t, "pracas_pedagio.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("antt.pracas-pedagio"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 5 {
		t.Fatalf("rowCount: got %d want 5", rowCount)
	}
}

func readANTTTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "antt", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
