package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestDNITSNVRodoviasFederaisGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readDNITIngestTestdata(t, "snv_rodovias_federais.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("dnit.snv-rodovias-federais"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func TestIPEASeriesMacroRegionaisGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIPEAIngestTestdata(t, "series_macro_regionais.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ipea.series-macro-regionais"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 56 {
		t.Fatalf("rowCount: got %d want 56", rowCount)
	}
}

func TestIBGEPEVSProducaoVegetalGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGEIngestTestdata(t, "pevs_producao_vegetal.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.pevs-producao-vegetal"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 14 {
		t.Fatalf("rowCount: got %d want 14", rowCount)
	}
}

func TestIBGEPPMProducaoMunicipalGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGEIngestTestdata(t, "ppm_producao_municipal.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.ppm-producao-municipal"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 7 {
		t.Fatalf("rowCount: got %d want 7", rowCount)
	}
}

func TestANEELTarifasEnergiaGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readANEELIngestTestdata(t, "tarifas_energia.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("aneel.tarifas-energia"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func TestBNDESFinanciamentoAgroGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readBNDESIngestTestdata(t, "financiamento_agro.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("bndes.financiamento-agro"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func TestINMETSequiaMonitorGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readINMETIngestTestdata(t, "sequia_monitor.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("inmet.sequia-monitor"),
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

func readDNITIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "dnit", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readIPEAIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "ipea", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readIBGEIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "ibge", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readANEELIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "aneel", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readBNDESIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "bndes", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readINMETIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "inmet", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
