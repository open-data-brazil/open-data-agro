package ingest

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestMDICComexImportGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readMDICTestdata(t, "comex_import.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mdic.comex-importacao-ncm-mes"),
		Format:    catalog.FormatJSON,
		ComexFlow: "import",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 2 {
		t.Fatalf("rowCount: got %d want 2", rowCount)
	}
}

func TestMDICComexExportUFGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readMDICTestdata(t, "comex_export_uf.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mdic.comex-exportacao-uf-ncm"),
		Format:    catalog.FormatJSON,
		ComexFlow: "export",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func TestMDICComexImportDieselGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readMDICTestdata(t, "comex_import_diesel.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mdic.comex-importacao-diesel-ncm"),
		Format:    catalog.FormatJSON,
		ComexFlow: "import",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 1 {
		t.Fatalf("rowCount: got %d want 1", rowCount)
	}
}
