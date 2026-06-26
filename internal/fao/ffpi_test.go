package fao

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenFFPIGoldenVector(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("testdata", "food_price_index.sample.json"))
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fao.food-price-index"),
	}
	headers, rows, err := FlattenFFPI(entry, raw)
	if err != nil {
		t.Fatalf("FlattenFFPI: %v", err)
	}
	if len(headers) != 5 || len(rows) != 12 {
		t.Fatalf("headers=%d rows=%d", len(headers), len(rows))
	}
}

func TestParseFFPICSVGoldenVector(t *testing.T) {
	t.Parallel()

	raw := []byte(`MONTHLY FOOD PRICE INDICES (2002-2004=100),,,,,,
,,,,,,
Date,Food Price Index,Meat Price Index,Dairy Price Index,Cereals Price Index,Oils Price Index,Sugar Price Index
Jan-24,118.0,110.0,120.0,110.5,105.0,130.0
`)
	entry := catalog.RegistryEntry{PeriodStart: 2020}
	rows, err := parseFFPICSV(raw, entry)
	if err != nil {
		t.Fatalf("parseFFPICSV: %v", err)
	}
	if len(rows) != 6 {
		t.Fatalf("rows: got %d want 6", len(rows))
	}
}
