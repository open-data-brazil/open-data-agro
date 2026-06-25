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

func readCONABTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "conab", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
