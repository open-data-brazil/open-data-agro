package igc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenGOI(t *testing.T) {
	t.Parallel()

	raw := readGOITestdata(t, "goi_index.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("igc.goi-index"),
		Format:    catalog.FormatJSON,
	}

	headers, rows, err := FlattenGOI(entry, raw)
	if err != nil {
		t.Fatalf("FlattenGOI: %v", err)
	}
	if len(headers) != 6 {
		t.Fatalf("headers: got %d want 6", len(headers))
	}
	if len(rows) != 4 {
		t.Fatalf("rows: got %d want 4", len(rows))
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	url, err := ResolveURL(catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("igc.goi-index"),
	})
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" {
		t.Fatal("expected default GOI URL")
	}
}

func readGOITestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
