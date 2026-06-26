package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestWorldBankPinkSheetGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readWorldBankIngestTestdata(t, "pink_sheet.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("worldbank.pink-sheet-monthly"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 5 {
		t.Fatalf("rowCount: got %d want 5", rowCount)
	}
}

func readWorldBankIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "worldbank", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
