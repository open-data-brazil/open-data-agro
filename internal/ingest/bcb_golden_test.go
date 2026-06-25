package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestBCBSGSIPCAGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readBCBTestdata(t, "sgs_ipca.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("bcb.sgs-ipca"),
		Format:    catalog.FormatJSON,
		SGSCode:   433,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func readBCBTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "bcb", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
