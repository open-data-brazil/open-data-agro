package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestOECDFAOAgOutlookGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readOECDTestdata(t, "ag_outlook.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("oecd-fao.ag-outlook"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 1 {
		t.Fatalf("rowCount: got %d want >= 1", rowCount)
	}
}

func TestFAOFFPIGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readFAOTestdata(t, "food_price_index.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fao.food-price-index"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 12 {
		t.Fatalf("rowCount: got %d want 12", rowCount)
	}
}

func TestArgentinaMAGyPGranosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readArgentinaTestdata(t, "magyp_producion_granos.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("argentina.magyp-producion-granos"),
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

func readOECDTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "oecd", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readFAOTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "fao", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
