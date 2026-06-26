package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestUSDAWASDEGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readUSDATestdataGolden(t, "wasde.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("usda.wasde"),
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

func TestUNComtradeBulkGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readUNTestdataGolden(t, "comtrade_brazil.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("un.comtrade-bulk"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 2 {
		t.Fatalf("rowCount: got %d want 2", rowCount)
	}
}

func readUSDATestdataGolden(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "usda", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readUNTestdataGolden(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "un", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
