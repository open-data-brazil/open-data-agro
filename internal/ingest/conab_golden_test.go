package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestLevantamentoGraosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "LevantamentoGraos.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 4 {
		t.Fatalf("rowCount: got %d want >= 4", rowCount)
	}
}

func TestSerieHistoricaGraosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "SerieHistoricaGraos.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.serie-historica-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 4 {
		t.Fatalf("rowCount: got %d want >= 4", rowCount)
	}
}

func TestOfertaDemandaGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "OfertaDemanda.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.oferta-demanda"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 4 {
		t.Fatalf("rowCount: got %d want >= 4", rowCount)
	}
}

func TestEstoquesPublicosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "Estoques.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estoques-publicos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 4 {
		t.Fatalf("rowCount: got %d want >= 4", rowCount)
	}
}

func TestANPCombustiveisMediosGoldenVector(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "anp", "testdata", "combustiveis_precos_medios.sample.xlsx")
	entry := catalog.RegistryEntry{
		DatasetID:     catalog.MustParseDatasetID("anp.combustiveis-precos-medios-municipios"),
		Format:        catalog.FormatXLSX,
		XLSXSheet:     "MUNICIPIOS",
		XLSXHeaderRow: 9,
	}

	_, rowCount, err := ConvertToParquetFromFile(entry, path)
	if err != nil {
		t.Fatalf("ConvertToParquetFromFile: %v", err)
	}
	if rowCount < 100 {
		t.Fatalf("rowCount: got %d want >= 100", rowCount)
	}
}

func TestANPCombustiveisPostosGoldenVector(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "anp", "testdata", "combustiveis_precos_postos.sample.xlsx")
	entry := catalog.RegistryEntry{
		DatasetID:     catalog.MustParseDatasetID("anp.combustiveis-precos-postos"),
		Format:        catalog.FormatXLSX,
		XLSXSheet:     "POSTOS REVENDEDORES",
		XLSXHeaderRow: 9,
	}

	_, rowCount, err := ConvertToParquetFromFile(entry, path)
	if err != nil {
		t.Fatalf("ConvertToParquetFromFile: %v", err)
	}
	if rowCount < 100 {
		t.Fatalf("rowCount: got %d want >= 100", rowCount)
	}
}

func TestArmazenagemGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "ArmazensCadastrados.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.armazenagem"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 4 {
		t.Fatalf("rowCount: got %d want >= 4", rowCount)
	}
}

func TestCapacidadeEstaticaGoldenVector(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "conab", "testdata", "capacidade_estatica.sample.xls")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.serie-historica-capacidade-estatica"),
		Format:    catalog.FormatXLS,
	}

	_, rowCount, err := ConvertToParquetFromFile(entry, path)
	if err != nil {
		t.Fatalf("ConvertToParquetFromFile: %v", err)
	}
	if rowCount < 100 {
		t.Fatalf("rowCount: got %d want >= 100", rowCount)
	}
}

func TestAlimentaBrasilEntregasGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "PAA_Entregas.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.alimenta-brasil-entregas"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 4 {
		t.Fatalf("rowCount: got %d want >= 4", rowCount)
	}
}

func TestAlimentaBrasilPropostasGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "PAA_Propostas.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.alimenta-brasil-propostas"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 4 {
		t.Fatalf("rowCount: got %d want >= 4", rowCount)
	}
}

func readCONABTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "conab", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
