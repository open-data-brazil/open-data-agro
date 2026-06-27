package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestIBAMASISFOGOGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBAMATestdata(t, "sisfogo_roi.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibama.sisfogo-incendios"),
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

func TestIBAMALicencasGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBAMATestdata(t, "sislic_licencas.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibama.licencas-ambientais"),
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

func TestIBAMAAutosGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readIBAMATestdata(t, "auto_infracao_1977.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ibama.autos-infracao"),
		Format:    catalog.FormatCSV,
		Delimiter: ";",
	}
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 2 {
		t.Fatalf("rowCount: got %d want 2", rowCount)
	}
}

func TestANAPluviometriaGoldenVector(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ana.pluviometria-redes"),
		Format:    catalog.FormatJSON,
	}
	raw := []byte(`[{"station_code":"87017001","consistency_level":"1","data_type":"2","observed_at":"2024-06-01 00:00:00","daily_mean":"12.4"},{"station_code":"87017001","consistency_level":"1","data_type":"2","observed_at":"2024-06-02 00:00:00","daily_mean":"0.0"}]`)
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 2 {
		t.Fatalf("rowCount: got %d want 2", rowCount)
	}
}

func TestEmbrapaAgroAPIGoldenVector(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("embrapa.agroapi-agrofit"),
		Format:    catalog.FormatJSON,
	}
	raw := []byte(`[{"numero_registro":"35523","marca_comercial":"KBR-829M1-02","situacao":"TRUE","classe":"Agente Biológico de Controle","formulacao":"Nematóides vivos","ingrediente_ativo":"Heterorhabditis bacteriophora"}]`)
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 1 {
		t.Fatalf("rowCount: got %d want 1", rowCount)
	}
}

func TestTransportesMTRBITShapefileGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readTransportesTestdata(t, "baseferro_attributes.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("transportes.mtr-bit-malha-shapefile"),
		Format:    catalog.FormatCSV,
		Delimiter: ",",
	}
	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 5 {
		t.Fatalf("rowCount: got %d want 5", rowCount)
	}
}

func readIBAMATestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "ibama", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func readTransportesTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "transportes", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
