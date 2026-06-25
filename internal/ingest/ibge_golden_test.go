package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestIBGEMunicipiosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "municipios.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.localidades-municipios"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 1 {
		t.Fatalf("rowCount: got %d want 1", rowCount)
	}
}

func TestIBGEUFsGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "ufs.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.localidades-ufs"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 27 {
		t.Fatalf("rowCount: got %d want 27", rowCount)
	}
}

func TestIBGERegioesGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "regioes.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.localidades-regioes"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 5 {
		t.Fatalf("rowCount: got %d want >= 5", rowCount)
	}
}

func readIBGETestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "ibge", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
