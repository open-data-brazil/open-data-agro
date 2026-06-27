package jrc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenMARSCropYieldGoldenVector(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("testdata", "mars_crop_yield.sample.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("jrc.mars-crop-yield"),
	}
	_, rows, err := FlattenMARSCropYield(entry, raw)
	if err != nil {
		t.Fatalf("FlattenMARSCropYield: %v", err)
	}
	if len(rows) < 1 {
		t.Fatalf("rows: got %d want >= 1", len(rows))
	}
}
