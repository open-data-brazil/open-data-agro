package ingest

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/extrame/xls"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func convertXLSToParquet(entry catalog.RegistryEntry, raw []byte) ([]byte, int, error) {
	xl, err := xls.OpenReader(bytes.NewReader(raw), "utf-8")
	if err != nil {
		return nil, 0, fmt.Errorf("open workbook: %w", err)
	}
	return convertXLSSheet(entry, xl)
}

func convertXLSFileToParquet(entry catalog.RegistryEntry, path string) ([]byte, int, error) {
	xl, err := xls.Open(path, "utf-8")
	if err != nil {
		return nil, 0, fmt.Errorf("open workbook: %w", err)
	}
	return convertXLSSheet(entry, xl)
}

func convertXLSSheet(entry catalog.RegistryEntry, xl *xls.WorkBook) ([]byte, int, error) {
	sheetIndex := 0
	if entry.XLSXSheet != "" {
		idx, err := strconv.Atoi(entry.XLSXSheet)
		if err != nil {
			return nil, 0, fmt.Errorf("xlsx_sheet must be numeric for xls datasets: %w", err)
		}
		sheetIndex = idx
	}

	sheet := xl.GetSheet(sheetIndex)
	if sheet == nil {
		return nil, 0, fmt.Errorf("workbook has no sheet at index %d", sheetIndex)
	}

	headerRow := entry.XLSXHeaderRow
	if headerRow < 0 || headerRow > int(sheet.MaxRow) {
		return nil, 0, fmt.Errorf("xlsx_header_row %d out of range", headerRow)
	}

	headers := normalizeHeaders(readXLSRow(sheet, headerRow))
	if len(headers) == 0 {
		return nil, 0, fmt.Errorf("empty header row %d", headerRow)
	}

	var rows [][]string
	for i := headerRow + 1; i <= int(sheet.MaxRow); i++ {
		record := readXLSRow(sheet, i)
		if rowIsEmpty(record) {
			continue
		}
		rows = append(rows, padRecord(record, len(headers)))
	}

	return writeStringTable(headers, rows)
}

func readXLSRow(sheet *xls.WorkSheet, rowIndex int) []string {
	row := sheet.Row(rowIndex)
	if row == nil {
		return nil
	}
	lastCol := row.LastCol()
	if lastCol <= 0 {
		return nil
	}
	out := make([]string, lastCol)
	for col := 0; col < lastCol; col++ {
		out[col] = row.Col(col)
	}
	return out
}
