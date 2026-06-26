package fao

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenProducaoAgro(t *testing.T) {
	t.Parallel()

	raw := readFAOTestdata(t, "producao_agro.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fao.producao-agro"),
	}

	headers, rows, err := Flatten(entry, raw)
	if err != nil {
		t.Fatalf("Flatten: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}
	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["commodity_slug"]]; got != "milho" {
		t.Fatalf("commodity_slug: got %q", got)
	}
	if got := rows[0][idx["element_code"]]; got != "5510" {
		t.Fatalf("element_code: got %q", got)
	}
}

func TestFlattenComercioAgro(t *testing.T) {
	t.Parallel()

	raw := readFAOTestdata(t, "comercio_agro.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fao.comercio-agro"),
	}

	headers, rows, err := Flatten(entry, raw)
	if err != nil {
		t.Fatalf("Flatten: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}
	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[1][idx["element_code"]]; got != "5911" {
		t.Fatalf("import element: got %q want 5911", got)
	}
}
