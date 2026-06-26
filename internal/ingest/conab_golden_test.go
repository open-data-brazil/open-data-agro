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

func TestPrecosMinimosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "PrecoMinimo.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.precos-minimos"),
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

func TestPrecosSemanalUFGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "PrecosSemanalUF.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.precos-agropecuarios-semanal-uf"),
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

func TestPrecosSemanalMunicipioGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "PrecosSemanalMunicipio.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.precos-agropecuarios-semanal-municipio"),
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

func TestPrecosMensalUFGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "PrecosMensalUF.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.precos-agropecuarios-mensal-uf"),
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

func TestPrecosMensalMunicipioGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "PrecosMensalMunicipio.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.precos-agropecuarios-mensal-municipio"),
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

func TestProhortDiarioGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "ProhortDiario.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.prohort-diario"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 5 {
		t.Fatalf("rowCount: got %d want 5", rowCount)
	}
}

func TestProhortMensalGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "ProhortMensal.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.prohort-mensal"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 5 {
		t.Fatalf("rowCount: got %d want 5", rowCount)
	}
}

func TestOperacoesComercializacaoGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "Leilao.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.operacoes-comercializacao"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 5 {
		t.Fatalf("rowCount: got %d want 5", rowCount)
	}
}

func TestVendasBalcaoGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "VendaBalcao.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.vendas-balcao"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 5 {
		t.Fatalf("rowCount: got %d want 5", rowCount)
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

func TestFreteGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCONABTestdata(t, "Frete.sample.txt")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.frete"),
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
