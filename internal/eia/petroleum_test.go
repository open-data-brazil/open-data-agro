package eia

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenPetroleumGolden(t *testing.T) {
	t.Parallel()

	raw := readEIATestdata(t, "petroleum_prices.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("eia.petroleum-prices"),
	}

	headers, rows, err := FlattenPetroleum(entry, raw)
	if err != nil {
		t.Fatalf("FlattenPetroleum: %v", err)
	}
	if len(headers) != 7 {
		t.Fatalf("headers: got %d want 7", len(headers))
	}
	if len(rows) != 4 {
		t.Fatalf("rows: got %d want 4", len(rows))
	}
	if rows[0][2] != "wti_spot" {
		t.Fatalf("commodity_slug: got %q", rows[0][2])
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("eia.petroleum-prices"),
		SourceURL: "https://api.eia.gov/v2/petroleum/pri/spt/data",
	}
	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if !strings.Contains(url, "api.eia.gov") {
		t.Fatalf("url: got %q", url)
	}
}

func TestBuildSeriesURL(t *testing.T) {
	t.Parallel()

	url := buildSeriesURL("PET.RWTC.D", "test-key", "2010-01-01", "2024-01-01", 0)
	if !strings.Contains(url, "/seriesid/PET.RWTC.D") {
		t.Fatalf("missing series path: %q", url)
	}
	if !strings.Contains(url, "api_key=test-key") {
		t.Fatalf("missing api key: %q", url)
	}
}

func readEIATestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
