package ingest

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestIBGEPPMEfetivoRebanhosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGEIngestTestdata(t, "ppm_herd.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.ppm-efetivo-rebanhos"),
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

func TestIBGEPAMPrecosProdutorGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGEIngestTestdata(t, "pam_area_quantidade.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.pam-precos-produtor"),
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

func TestIBGECensoAgroAreaUsoSoloGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGEIngestTestdata(t, "censo_agro_estabelecimentos.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.censo-agro-area-uso-solo"),
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

func TestIBGEPNADRuralRendaOcupacaoGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBGEIngestTestdata(t, "pnad_continua_rural.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibge.pnad-rural-renda-ocupacao"),
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
