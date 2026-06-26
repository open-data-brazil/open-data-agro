package b3

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestFlattenFuturoSojaGolden(t *testing.T) {
	t.Parallel()

	raw := readB3Testdata(t, "futuro_soja.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID:         catalog.MustParseDatasetID("b3.futuro-soja"),
		B3CommodityPrefix: "SOY",
	}

	headers, rows, err := FlattenFuturo(entry, raw)
	if err != nil {
		t.Fatalf("FlattenFuturo: %v", err)
	}
	if len(headers) != 8 {
		t.Fatalf("headers: got %d want 8", len(headers))
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}
	if rows[0][1] != "SOYH26" {
		t.Fatalf("symbol: got %q", rows[0][1])
	}
}

func TestParseSPRDXMLSample(t *testing.T) {
	t.Parallel()

	raw := readB3Testdata(t, "sprd_soja.sample.xml")
	rows, err := parseSPRDXML(raw, "SOY")
	if err != nil {
		t.Fatalf("parseSPRDXML: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows: got %d want 1", len(rows))
	}
	if rows[0].Symbol != "SOYH26" {
		t.Fatalf("symbol: got %q", rows[0].Symbol)
	}
}

func readB3Testdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
