package ingest

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func convertXLSXToParquet(raw []byte) ([]byte, int, error) {
	book, err := excelize.OpenReader(bytes.NewReader(raw))
	if err != nil {
		return nil, 0, fmt.Errorf("open workbook: %w", err)
	}
	defer func() { _ = book.Close() }()
	return convertWorkbookSheet(book)
}

func convertXLSXFileToParquet(path string) ([]byte, int, error) {
	book, err := excelize.OpenFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("open workbook: %w", err)
	}
	defer func() { _ = book.Close() }()
	return convertWorkbookSheet(book)
}

func convertWorkbookSheet(book *excelize.File) ([]byte, int, error) {
	sheet := book.GetSheetName(0)
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

	headers := normalizeHeaders(table[0])
	var rows [][]string
	for _, record := range table[1:] {
		rows = append(rows, padRecord(record, len(headers)))
	}

	return writeStringTable(headers, rows)
}
