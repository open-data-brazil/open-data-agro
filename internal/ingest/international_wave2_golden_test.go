package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestUSDAGATSTradeGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readUSDATestdata(t, "gats_trade.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("usda.gats-trade"),
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

func TestEurostatAgPricesGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readEurostatTestdata(t, "ag_prices.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("eurostat.ag-prices"),
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

func TestArgentinaBCRACambioGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readArgentinaTestdata(t, "bcra_cambio.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("argentina.bcra-cambio"),
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

func readUSDATestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "usda", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readEurostatTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "eurostat", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readArgentinaTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "argentina", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
