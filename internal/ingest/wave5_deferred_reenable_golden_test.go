package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// Golden vectors for deferred datasets — parser readiness when live URLs unblock.

func TestDeferredReenableANTAQMovimentacaoCargaGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readDeferredAgencyTestdata(t, "antaq", "movimentacao_carga_portuaria.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("antaq.movimentacao-carga-portuaria"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 1 {
		t.Fatalf("rowCount: got %d want >= 1", rowCount)
	}
}

func TestDeferredReenableFAOComercioAgroGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readFAOIngestTestdata(t, "comercio_agro.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fao.comercio-agro"),
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

func TestDeferredReenableUSDAGATSTradeGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readDeferredUSDATestdata(t, "gats_trade.sample.json")
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

func TestDeferredReenableUSDAPSDSojaGoldenVector(t *testing.T) {
	t.Parallel()
	runDeferredPSDGolden(t, "usda.psd-soja")
}

func TestDeferredReenableUSDAPSDMilhoGoldenVector(t *testing.T) {
	t.Parallel()
	runDeferredPSDGolden(t, "usda.psd-milho")
}

func TestDeferredReenableUSDAPSDTrigoGoldenVector(t *testing.T) {
	t.Parallel()
	runDeferredPSDGolden(t, "usda.psd-trigo")
}

func TestDeferredReenableWTOITSTradeGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readDeferredAgencyTestdata(t, "wto", "its_trade_statistics.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("wto.its-trade-statistics"),
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

func TestDeferredReenableMexicoSIAPProduccionGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readDeferredAgencyTestdata(t, "mexico", "siap_produccion_agricola.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mexico.siap-produccion-agricola"),
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

func TestDeferredReenableNOAAGPCCPrecipitationGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readDeferredAgencyTestdata(t, "noaa", "gpcc_precipitation.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("noaa.gpcc-precipitation"),
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

func runDeferredPSDGolden(t *testing.T, datasetID string) {
	t.Helper()

	raw := readDeferredUSDATestdata(t, "psd_soja.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID(datasetID),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 4 {
		t.Fatalf("rowCount: got %d want 4", rowCount)
	}
}

func readDeferredAgencyTestdata(t *testing.T, agency, name string) []byte {
	t.Helper()
	path := filepath.Join("..", agency, "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readDeferredUSDATestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "usda", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
