package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestCFTCCOTAgriculturalGoldenVector(t *testing.T) {
	t.Parallel()
	raw := readAgencyTestdata(t, "cftc", "cot_agricultural.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("cftc.cot-agricultural-futures"),
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

func TestJRCMARSCropYieldGoldenVector(t *testing.T) {
	t.Parallel()
	raw := readAgencyTestdata(t, "jrc", "mars_crop_yield.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("jrc.mars-crop-yield"),
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

func TestFAOGIEWSCropProspectsGoldenVector(t *testing.T) {
	t.Parallel()
	raw := readFAOTestdata(t, "giews_crop_prospects.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fao.giews-crop-prospects"),
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

func TestFAOAMISMarketMonitorGoldenVector(t *testing.T) {
	t.Parallel()
	raw := readFAOTestdata(t, "amis_market_monitor.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fao.amis-market-monitor"),
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

func TestSAGISGrainSupplyGoldenVector(t *testing.T) {
	t.Parallel()
	raw := readAgencyTestdata(t, "sagis", "grain_supply_statistics.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("sagis.grain-supply-statistics"),
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

func TestJapanMAFFAgTradeGoldenVector(t *testing.T) {
	t.Parallel()
	raw := readAgencyTestdata(t, "japan", "maff_ag_trade.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("japan.maff-ag-trade"),
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

func TestFREDCommodityIndexesGoldenVector(t *testing.T) {
	t.Parallel()
	raw := readAgencyTestdata(t, "fred", "commodity_indexes.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fred.commodity-indexes"),
		Format:    catalog.FormatJSON,
	}
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount >= 1 {
		return
	}
	t.Fatalf("rowCount: got %d want >= 1", rowCount)
}

func TestNASAPOWERAgroGoldenVector(t *testing.T) {
	t.Parallel()
	raw := readAgencyTestdata(t, "nasa", "power_agroclimatology.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("nasa.power-agroclimatology"),
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

func TestCopernicusERA5AgroclimateGoldenVector(t *testing.T) {
	t.Parallel()
	raw := readAgencyTestdata(t, "copernicus", "era5_agroclimate.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("copernicus.era5-agroclimate"),
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

func readAgencyTestdata(t *testing.T, agency, name string) []byte {
	t.Helper()
	path := filepath.Join("..", agency, "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
