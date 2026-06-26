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
