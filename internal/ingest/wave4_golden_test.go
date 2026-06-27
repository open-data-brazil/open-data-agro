package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestIBGECensoAgroEstabelecimentosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGEIngestTestdata(t, "censo_agro_estabelecimentos.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.censo-agro-estabelecimentos"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 55 {
		t.Fatalf("rowCount: got %d want 55", rowCount)
	}
}

func TestIBGEPNADContinuaRuralGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGEIngestTestdata(t, "pnad_continua_rural.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.pnad-continua-rural"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 33 {
		t.Fatalf("rowCount: got %d want 33", rowCount)
	}
}

func TestSUFRAMAComercioMercadoriasZFMGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readSUFRAMAIngestTestdata(t, "comercio_mercadorias_zfm.sample.xlsx")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("suframa.comercio-mercadorias-zfm"),
		Format:    catalog.FormatXLSX,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 1 {
		t.Fatalf("rowCount: got %d want >= 1", rowCount)
	}
}

func TestTransportesMTRBITMalhaRodoviariaGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readTransportesIngestTestdata(t, "mtr_bit_malha_rodoviaria.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("transportes.mtr-bit-malha-rodoviaria"),
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

func TestMAPASIFAbateEstatisticasGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readMAPAIngestTestdata(t, "sif_abate_estatisticas.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mapa.sif-abate-estatisticas"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 14 {
		t.Fatalf("rowCount: got %d want 14", rowCount)
	}
}

func TestONSCargaEnergeticaGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readONSIngestTestdata(t, "carga_energetica.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ons.carga-energetica"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 3 {
		t.Fatalf("rowCount: got %d want >= 3", rowCount)
	}
}

func TestINPEDETERAlertasDesmatamentoGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readINPEIngestTestdata(t, "deter_alertas_desmatamento.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("inpe.deter-alertas-desmatamento"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 5 {
		t.Fatalf("rowCount: got %d want 5", rowCount)
	}
}

func TestDNITCondicoesConservacaoRodoviasGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readDNITIngestTestdata(t, "condicoes_conservacao_rodovias.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("dnit.condicoes-conservacao-rodovias"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 14 {
		t.Fatalf("rowCount: got %d want 14", rowCount)
	}
}

func readSUFRAMAIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "suframa", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readTransportesIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "transportes", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readONSIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "ons", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readINPEIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "inpe", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readMAPAIngestTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "mapa", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
