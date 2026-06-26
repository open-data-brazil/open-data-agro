package ingest

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFAOProducaoAgroGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readFAOIngestTestdata(t, "producao_agro.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fao.producao-agro"),
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

func TestFAOComercioAgroGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readFAOIngestTestdata(t, "comercio_agro.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fao.comercio-agro"),
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

func TestWorldBankAgIndicesGoldenVector(t *testing.T) {
	t.Parallel()

	raw := readWorldBankIngestTestdata(t, "ag_indices.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("worldbank.ag-indices"),
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
