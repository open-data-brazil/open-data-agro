package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestNOAAENSOGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readNOAAIngestTestdata(t, "enso.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("noaa.enso-indices"),
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

func TestNOAAGlobalTempGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readNOAAIngestTestdata(t, "global_temp.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("noaa.global-temp-anomaly"),
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

func readNOAAIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "noaa", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
