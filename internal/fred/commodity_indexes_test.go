package fred

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenCommodityIndexesGoldenVector(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("testdata", "commodity_indexes.sample.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fred.commodity-indexes"),
	}
	_, rows, err := FlattenCommodityIndexes(entry, raw)
	if err != nil {
		t.Fatalf("FlattenCommodityIndexes: %v", err)
	}
	if len(rows) < 1 {
		t.Fatalf("rows: got %d want >= 1", len(rows))
	}
}
