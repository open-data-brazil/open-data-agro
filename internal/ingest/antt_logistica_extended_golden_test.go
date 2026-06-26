package ingest

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestANTTVolumeTrafegoPedagioGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readANTTTestdata(t, "volume_trafego_pedagio.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("antt.volume-trafego-pedagio"),
		Format:    catalog.FormatCSV,
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

func TestANTTReceitaPorPracaGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readANTTTestdata(t, "receita_por_praca.sample.csv")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("antt.receita-por-praca"),
		Format:    catalog.FormatCSV,
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
