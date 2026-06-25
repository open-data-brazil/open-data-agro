package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestCEPASojaParanaguaGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readCEPATestdata(t, "soja_paranagua_historico.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("cepea.soja-paranagua"),
		Format:    catalog.FormatJSON,
		License:   "CC BY-NC 4.0",
		FonteTipo: "referencia_mercado",
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 5 {
		t.Fatalf("rowCount: got %d want 5", rowCount)
	}
}

func TestCEPABronzeMetadataLicense(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID:   catalog.MustParseDatasetID("cepea.soja-paranagua"),
		PortalLabel: "Soja Paranaguá",
		License:     "CC BY-NC 4.0",
		FonteTipo:   "referencia_mercado",
	}
	meta := NewBronzeMetadata(entry, Fingerprint{}, "https://example.test", "local")
	if meta.Licenca != "CC BY-NC 4.0" {
		t.Fatalf("licenca: got %q", meta.Licenca)
	}
	if meta.Agencia != "CEPEA" {
		t.Fatalf("agencia: got %q", meta.Agencia)
	}
	if meta.FonteTipo != "referencia_mercado" {
		t.Fatalf("fonte_tipo: got %q", meta.FonteTipo)
	}
}

func readCEPATestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "cepea", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
