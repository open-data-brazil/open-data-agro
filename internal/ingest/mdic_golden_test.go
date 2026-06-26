package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestMDICComexExportGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readMDICTestdata(t, "comex_export.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mdic.comex-exportacao-ncm-mes"),
		Format:    catalog.FormatJSON,
		ComexFlow: "export",
		ComexNCMs: []string{"12019000", "10059000"},
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func readMDICTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "mdic", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
