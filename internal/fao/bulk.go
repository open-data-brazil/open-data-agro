package fao

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type bulkRow struct {
	AreaCode      string `json:"area_code"`
	AreaName      string `json:"area_name"`
	ItemCode      string `json:"item_code"`
	ItemName      string `json:"item_name"`
	ElementCode   string `json:"element_code"`
	ElementName   string `json:"element_name"`
	Year          string `json:"year"`
	MonthsCode    string `json:"months_code,omitempty"`
	Months        string `json:"months,omitempty"`
	Unit          string `json:"unit"`
	Value         string `json:"value"`
	Flag          string `json:"flag"`
	CommoditySlug string `json:"commodity_slug"`
}

func parseFAOBulkZip(raw []byte, csvName string, itemCodes, elementCodes map[string]struct{}, startYear, endYear int, withMonths bool) ([]bulkRow, error) {
	if strings.TrimSpace(csvName) == "" {
		return nil, fmt.Errorf("empty fao bulk csv name")
	}

	reader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return nil, fmt.Errorf("open fao bulk zip: %w", err)
	}

	var csvFile *zip.File
	for _, file := range reader.File {
		if file.Name == csvName {
			csvFile = file
			break
		}
	}
	if csvFile == nil {
		return nil, fmt.Errorf("csv %q not found in fao bulk zip", csvName)
	}

	rc, err := csvFile.Open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = rc.Close() }()

	return parseFAOCSV(rc, itemCodes, elementCodes, startYear, endYear, withMonths)
}

func parseFAOCSV(r io.Reader, itemCodes, elementCodes map[string]struct{}, startYear, endYear int, withMonths bool) ([]bulkRow, error) {
	csvReader := csv.NewReader(r)
	csvReader.ReuseRecord = true

	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("read fao bulk header: %w", err)
	}
	col := indexColumns(header)

	var rows []bulkRow
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read fao bulk csv: %w", err)
		}

		itemCode := fieldAt(record, col, "Item Code")
		elementCode := fieldAt(record, col, "Element Code")
		if !containsCode(itemCodes, itemCode) || !containsCode(elementCodes, elementCode) {
			continue
		}

		year, err := strconv.Atoi(strings.TrimSpace(fieldAt(record, col, "Year")))
		if err != nil || year < startYear || year > endYear {
			continue
		}

		row := bulkRow{
			AreaCode:      fieldAt(record, col, "Area Code"),
			AreaName:      fieldAt(record, col, "Area"),
			ItemCode:      itemCode,
			ItemName:      fieldAt(record, col, "Item"),
			ElementCode:   elementCode,
			ElementName:   fieldAt(record, col, "Element"),
			Year:          strconv.Itoa(year),
			Unit:          fieldAt(record, col, "Unit"),
			Value:         fieldAt(record, col, "Value"),
			Flag:          fieldAt(record, col, "Flag"),
			CommoditySlug: ItemSlug[itemCode],
		}
		if withMonths {
			row.MonthsCode = fieldAt(record, col, "Months Code")
			row.Months = fieldAt(record, col, "Months")
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no fao bulk rows after filter")
	}
	return rows, nil
}

func indexColumns(header []string) map[string]int {
	out := make(map[string]int, len(header))
	for i, name := range header {
		out[strings.TrimSpace(name)] = i
	}
	return out
}

func fieldAt(record []string, col map[string]int, name string) string {
	idx, ok := col[name]
	if !ok || idx >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[idx])
}

func containsCode(set map[string]struct{}, code string) bool {
	_, ok := set[strings.TrimSpace(code)]
	return ok
}

func codeSet(codes []string) map[string]struct{} {
	out := make(map[string]struct{}, len(codes))
	for _, code := range codes {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		out[code] = struct{}{}
	}
	return out
}

func bulkRowKey(row bulkRow) string {
	parts := []string{row.AreaCode, row.ItemCode, row.ElementCode, row.Year}
	if row.MonthsCode != "" {
		parts = append(parts, row.MonthsCode)
	}
	return strings.Join(parts, "|")
}
