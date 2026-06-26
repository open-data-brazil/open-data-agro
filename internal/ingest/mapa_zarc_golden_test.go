package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestMAPAZARCTabuaRiscoGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readMAPATestdata(t, "zarc_tabua_risco.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mapa.zarc-tabua-risco"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func readMAPATestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "mapa", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
