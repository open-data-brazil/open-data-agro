package ingest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestINMETEstacoesAutomaticasGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readINMETTestdata(t, "estacoes_automaticas.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("inmet.estacoes-automaticas"),
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

func TestINMETBDMEPDiarioGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readINMETTestdata(t, "bdmep_daily_long.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("inmet.bdmep-diario"),
		Format:    catalog.FormatCSV,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 3 {
		t.Fatalf("rowCount: got %d want >= 3", rowCount)
	}
}

func TestINMETPacoteAnualAutomaticasGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readINMETTestdata(t, "bdmep_daily_long.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("inmet.pacote-anual-automaticas"),
		Format:    catalog.FormatCSV,
	}

	_, rowCount, err := ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount < 3 {
		t.Fatalf("rowCount: got %d want >= 3", rowCount)
	}
}

func readINMETTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "inmet", "testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
