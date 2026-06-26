package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestMAPAAgrofitProdutosFormuladosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readMAPATestdata(t, "agrofit_produtos_formulados.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mapa.agrofit-produtos-formulados"),
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

func TestIGCGOIIndexGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIGCTestdata(t, "goi_index.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("igc.goi-index"),
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

func readIGCTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "igc", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
