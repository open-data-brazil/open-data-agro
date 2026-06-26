package fao

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenPricesGolden(t *testing.T) {
	t.Parallel()

	raw := readFAOTestdata(t, "prices_agro.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("fao.prices-agro"),
	}

	headers, rows, err := FlattenPrices(entry, raw)
	if err != nil {
		t.Fatalf("FlattenPrices: %v", err)
	}
	if len(headers) != 13 {
		t.Fatalf("headers: got %d want 13", len(headers))
	}
	if len(rows) != 4 {
		t.Fatalf("rows: got %d want 4", len(rows))
	}
	if rows[0][4] != "soja" {
		t.Fatalf("commodity_slug: got %q", rows[0][4])
	}
}

func TestParsePricesCSVSample(t *testing.T) {
	t.Parallel()

	raw := readFAOTestdata(t, "prices_agro.sample.csv")
	items := codeSet(defaultItemCodes)
	elements := codeSet(defaultElementCodes)
	rows, err := parseFAOCSV(strings.NewReader(string(raw)), items, elements, 2010, 2030, true)
	if err != nil {
		t.Fatalf("parseFAOCSV: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}
	if rows[0].CommoditySlug != "soja" {
		t.Fatalf("slug: got %q", rows[0].CommoditySlug)
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID:  catalog.MustParseDatasetID("fao.prices-agro"),
		FAOBulkURL: defaultPricesBulkURL,
	}
	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if !strings.Contains(url, "bulks-faostat.fao.org") {
		t.Fatalf("url: got %q", url)
	}
}

func readFAOTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
