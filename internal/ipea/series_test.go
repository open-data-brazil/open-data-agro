package ipea

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenSeriesGolden(t *testing.T) {
	t.Parallel()

	raw := readIPEATestdata(t, "series_macro_regionais.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("ipea.series-macro-regionais"),
	}

	headers, rows, err := FlattenSeries(entry, raw)
	if err != nil {
		t.Fatalf("FlattenSeries: %v", err)
	}
	if len(rows) < 2 {
		t.Fatalf("rows: got %d want >= 2", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["series_code"]]; got == "" {
		t.Fatalf("series_code should not be empty")
	}
	if got := rows[0][idx["refdate"]]; len(got) != 10 {
		t.Fatalf("refdate: got %q want YYYY-MM-DD", got)
	}
}

func readIPEATestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
