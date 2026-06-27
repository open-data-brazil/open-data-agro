package ingest

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/anp"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/xuri/excelize/v2"
)

func convertXLSXToParquet(entry catalog.RegistryEntry, raw []byte) ([]byte, int, error) {
	book, err := excelize.OpenReader(bytes.NewReader(raw))
	if err != nil {
		return nil, 0, fmt.Errorf("open workbook: %w", err)
	}
	defer func() { _ = book.Close() }()
	return convertWorkbookSheet(entry, book)
}

func convertXLSXFileToParquet(entry catalog.RegistryEntry, path string) ([]byte, int, error) {
	book, err := excelize.OpenFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("open workbook: %w", err)
	}
	defer func() { _ = book.Close() }()
	return convertWorkbookSheet(entry, book)
}

func convertWorkbookSheet(entry catalog.RegistryEntry, book *excelize.File) ([]byte, int, error) {
	sheet := strings.TrimSpace(entry.XLSXSheet)
	if sheet == "" {
		sheet = book.GetSheetName(0)
	}
	if sheet == "" {
		return nil, 0, fmt.Errorf("workbook has no sheets")
	}

	table, err := book.GetRows(sheet)
	if err != nil {
		return nil, 0, fmt.Errorf("read sheet %q: %w", sheet, err)
	}
	if len(table) == 0 {
		return nil, 0, fmt.Errorf("sheet %q is empty", sheet)
	}

	headerRow := entry.XLSXHeaderRow
	if headerRow < 0 || headerRow >= len(table) {
		return nil, 0, fmt.Errorf("xlsx_header_row %d out of range for sheet %q", headerRow, sheet)
	}

	headers := normalizeHeaders(table[headerRow])
	if len(headers) == 0 {
		return nil, 0, fmt.Errorf("empty header row %d in sheet %q", headerRow, sheet)
	}

	var rows [][]string
	for _, record := range table[headerRow+1:] {
		if rowIsEmpty(record) {
			continue
		}
		rows = append(rows, padRecord(record, len(headers)))
	}

	if entry.DatasetID.String() == "anp.etanol-precos" {
		headers, rows = anp.FilterEthanolPrecos(headers, rows)
		if len(rows) == 0 {
			return nil, 0, fmt.Errorf("no ethanol rows in LPC sheet %q", sheet)
		}
	}

	return writeStringTable(headers, rows)
}

func rowIsEmpty(record []string) bool {
	for _, cell := range record {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
