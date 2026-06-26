package usda

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestParseWASDEXMLSample(t *testing.T) {
	t.Parallel()

	raw := readWASDETestdata(t, "wasde.sample.xml")
	rows, err := parseWASDEXML(raw, "2026-06")
	if err != nil {
		t.Fatalf("parseWASDEXML: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("expected rows from sample xml")
	}
	if rows[0].Commodity == "" || rows[0].Attribute == "" {
		t.Fatalf("incomplete row: %+v", rows[0])
	}
}

func TestFlattenWASDE(t *testing.T) {
	t.Parallel()

	payload := []byte(`[{"report_month":"2026-06","commodity":"Wheat","market_year":"2025/26","attribute":"Output","value":"800.1","unit":"million metric tons"}]`)
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("usda.wasde"),
	}

	headers, rows, err := FlattenWASDE(entry, payload)
	if err != nil {
		t.Fatalf("FlattenWASDE: %v", err)
	}
	if len(headers) != 6 {
		t.Fatalf("headers: got %d want 6", len(headers))
	}
	if len(rows) != 1 {
		t.Fatalf("rows: got %d want 1", len(rows))
	}
}

func TestWASDEResolveURL(t *testing.T) {
	t.Parallel()

	url, err := ResolveURL(catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("usda.wasde"),
	})
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" {
		t.Fatal("expected wasde index url")
	}
}

func readWASDETestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
