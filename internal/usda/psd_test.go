package usda

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenPSDSojaGolden(t *testing.T) {
	t.Parallel()

	raw := readUSDATestdata(t, "psd_soja.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID:        catalog.MustParseDatasetID("usda.psd-soja"),
		PSDCommodityCode: "2222000",
		PSDCommoditySlug: "soja",
	}

	headers, rows, err := FlattenPSD(entry, raw)
	if err != nil {
		t.Fatalf("FlattenPSD: %v", err)
	}
	if len(headers) != 13 {
		t.Fatalf("headers: got %d want 13", len(headers))
	}
	if len(rows) != 4 {
		t.Fatalf("rows: got %d want 4", len(rows))
	}
	if rows[0][3] != "BR" {
		t.Fatalf("country_code: got %q", rows[0][3])
	}
}

func TestParsePSDSoapSample(t *testing.T) {
	t.Parallel()

	raw := readUSDATestdata(t, "psd_soja.sample.xml")
	rows, err := parsePSDSoapResponse(raw)
	if err != nil {
		t.Fatalf("parsePSDSoapResponse: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows: got %d want 2", len(rows))
	}
	if rows[0].CountryCode != "BR" {
		t.Fatalf("country: got %q", rows[0].CountryCode)
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID:        catalog.MustParseDatasetID("usda.psd-soja"),
		PSDCommodityCode: "2222000",
	}
	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" {
		t.Fatal("empty url")
	}
}

func readUSDATestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
