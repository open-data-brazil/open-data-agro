package nasa

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenPOWERAgroGoldenVector(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("testdata", "power_agroclimatology.sample.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("nasa.power-agroclimatology"),
	}
	_, rows, err := FlattenPOWERAgro(entry, raw)
	if err != nil {
		t.Fatalf("FlattenPOWERAgro: %v", err)
	}
	if len(rows) < 1 {
		t.Fatalf("rows: got %d want >= 1", len(rows))
	}
}
