package worldbank

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/xuri/excelize/v2"
)

func TestFlattenPinkSheetGolden(t *testing.T) {
	t.Parallel()

	raw := readWorldBankTestdata(t, "pink_sheet.sample.json")
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("worldbank.pink-sheet-monthly"),
	}

	headers, rows, err := FlattenPinkSheet(entry, raw)
	if err != nil {
		t.Fatalf("FlattenPinkSheet: %v", err)
	}
	if len(headers) != 5 {
		t.Fatalf("headers: got %d want 5", len(headers))
	}
	if len(rows) != 5 {
		t.Fatalf("rows: got %d want 5", len(rows))
	}
	if rows[0][2] != "soja" {
		t.Fatalf("commodity_slug: got %q", rows[0][2])
	}
}

func TestParsePinkSheetXLSXSample(t *testing.T) {
	t.Parallel()

	raw := buildPinkSheetSampleXLSX(t)
	rows, err := parsePinkSheetXLSX(raw, defaultPinkSheetSheet, seriesNameSet(defaultSeriesNames), "2024-01", "2024-02")
	if err != nil {
		t.Fatalf("parsePinkSheetXLSX: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows: got %d want 2", len(rows))
	}
	if rows[0].CommoditySlug != "soja" {
		t.Fatalf("slug: got %q", rows[0].CommoditySlug)
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID:            catalog.MustParseDatasetID("worldbank.pink-sheet-monthly"),
		WorldBankPinkSheetURL: defaultPinkSheetURL,
	}
	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" {
		t.Fatal("empty url")
	}
}

func buildPinkSheetSampleXLSX(t *testing.T) []byte {
	t.Helper()

	book := excelize.NewFile()
	sheet := defaultPinkSheetSheet
	if err := book.SetSheetName("Sheet1", sheet); err != nil {
		t.Fatalf("rename sheet: %v", err)
	}

	_ = book.SetCellValue(sheet, "B5", "Soybeans")
	_ = book.SetCellValue(sheet, "B6", "($/mt)")
	_ = book.SetCellValue(sheet, "A7", "2024M01")
	_ = book.SetCellValue(sheet, "B7", 490.5)
	_ = book.SetCellValue(sheet, "A8", "2024M02")
	_ = book.SetCellValue(sheet, "B8", 500.1)

	var buf bytes.Buffer
	if err := book.Write(&buf); err != nil {
		t.Fatalf("write xlsx: %v", err)
	}
	return buf.Bytes()
}

func readWorldBankTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
