package worldbank

import (
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenAgIndices(t *testing.T) {
	t.Parallel()

	raw := readWorldBankTestdata(t, "ag_indices.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("worldbank.ag-indices"),
	}

	headers, rows, err := Flatten(entry, raw)
	if err != nil {
		t.Fatalf("Flatten: %v", err)
	}
	if len(rows) != 4 {
		t.Fatalf("rows: got %d want 4", len(rows))
	}
	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["commodity_slug"]]; got != "agriculture" {
		t.Fatalf("commodity_slug: got %q want agriculture", got)
	}
}
