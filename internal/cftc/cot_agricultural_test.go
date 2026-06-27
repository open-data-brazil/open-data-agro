package cftc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenCOTAgriculturalGoldenVector(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("testdata", "cot_agricultural.sample.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("cftc.cot-agricultural-futures"),
	}
	headers, rows, err := FlattenCOTAgricultural(entry, raw)
	if err != nil {
		t.Fatalf("FlattenCOTAgricultural: %v", err)
	}
	if len(headers) != 11 || len(rows) < 1 {
		t.Fatalf("headers=%d rows=%d", len(headers), len(rows))
	}
}
