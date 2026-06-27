package ingest

import (
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestBCBCimAgroCreditoRuralGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readBCBTestdata(t, "sgs_cim_agro.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("bcb.cim-agro-credito-rural"),
		Format:    catalog.FormatJSON,
		SGSCode:   21087,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 3 {
		t.Fatalf("rowCount: got %d want 3", rowCount)
	}
}

func TestBNDESDesembolsosLinhasAgroGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readBNDESIngestTestdata(t, "desembolsos_linhas_agro.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("bndes.desembolsos-linhas-agro"),
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

func TestANPEtanolPrecosGoldenVector(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "anp", "testdata", "combustiveis_precos_medios.sample.xlsx")
	entry := catalog.RegistryEntry{
		DatasetID:     catalog.MustParseDatasetID("anp.etanol-precos"),
		Format:        catalog.FormatXLSX,
		XLSXSheet:     "MUNICIPIOS",
		XLSXHeaderRow: 9,
	}

	_, rowCount, err := ConvertToParquetFromFile(entry, path)
	if err != nil {
		t.Fatalf("ConvertToParquetFromFile: %v", err)
	}
	if rowCount != 374 {
		t.Fatalf("rowCount: got %d want 374 (ethanol-only filter)", rowCount)
	}
}
