package un

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenComtrade(t *testing.T) {
	t.Parallel()

	raw := readUNTestdata(t, "comtrade_brazil.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("un.comtrade-bulk"),
	}

	headers, rows, err := FlattenComtrade(entry, raw)
	if err != nil {
		t.Fatalf("FlattenComtrade: %v", err)
	}
	if len(headers) != 13 {
		t.Fatalf("headers: got %d want 13", len(headers))
	}
	if len(rows) != 2 {
		t.Fatalf("rows: got %d want 2", len(rows))
	}
	if rows[0][8] != "soja" {
		t.Fatalf("commodity_slug: got %q", rows[0][8])
	}
}

func TestComtradeResolveURL(t *testing.T) {
	t.Parallel()

	url, err := ResolveURL(catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("un.comtrade-bulk"),
		SourceURL: "https://comtradeapi.un.org/data/v1/get/C/A/HS",
	})
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" {
		t.Fatal("expected comtrade url")
	}
}

func readUNTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
