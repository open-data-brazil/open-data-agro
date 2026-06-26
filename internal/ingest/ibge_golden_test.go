package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestIBGEMunicipiosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "municipios.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.localidades-municipios"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 1 {
		t.Fatalf("rowCount: got %d want 1", rowCount)
	}
}

func TestIBGEUFsGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "ufs.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.localidades-ufs"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 27 {
		t.Fatalf("rowCount: got %d want 27", rowCount)
	}
}

func TestIBGEPAMAreaQuantidadeGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "pam_area_quantidade.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.pam-area-quantidade"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 6 {
		t.Fatalf("rowCount: got %d want 6", rowCount)
	}
}

func TestIBGEPAMRendimentoValorGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "pam_rendimento_valor.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.pam-rendimento-valor"),
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

func TestIBGEPAMEstabelecimentosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "pam_estabelecimentos.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.pam-estabelecimentos"),
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

func TestIBGERegioesGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "regioes.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.localidades-regioes"),
		Format:    catalog.FormatJSON,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 5 {
		t.Fatalf("rowCount: got %d want >= 5", rowCount)
	}
}

func TestIBGEMesorregioesGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "mesorregioes.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.localidades-mesorregioes"),
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

func TestIBGEMicrorregioesGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "microrregioes.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.localidades-microrregioes"),
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

func TestIBGELSPAAreaProducaoGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "lspa_area_producao.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.lspa-area-producao"),
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

func readIBGETestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "ibge", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
