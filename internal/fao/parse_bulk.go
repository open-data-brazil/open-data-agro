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

const defaultPricesCSV = "Prices_E_All_Data_(Normalized).csv"

func parsePricesBulkZip(raw []byte, csvName string, itemCodes, elementCodes map[string]struct{}, startYear, endYear int) ([]priceRow, error) {
	if strings.TrimSpace(csvName) == "" {
		csvName = defaultPricesCSV
	}

	reader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return nil, fmt.Errorf("open fao prices zip: %w", err)
	}

	var csvFile *zip.File
	for _, file := range reader.File {
		if file.Name == csvName {
			csvFile = file
			break
		}
	}
	if csvFile == nil {
		return nil, fmt.Errorf("csv %q not found in fao prices zip", csvName)
	}

	rc, err := csvFile.Open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = rc.Close() }()

	return parsePricesCSV(rc, itemCodes, elementCodes, startYear, endYear)
}

func parsePricesCSV(r io.Reader, itemCodes, elementCodes map[string]struct{}, startYear, endYear int) ([]priceRow, error) {
	csvReader := csv.NewReader(r)
	csvReader.ReuseRecord = true

	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("read fao prices header: %w", err)
	}
	col := indexColumns(header)

	var rows []priceRow
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read fao prices csv: %w", err)
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

		row := priceRow{
			AreaCode:      fieldAt(record, col, "Area Code"),
			AreaName:      fieldAt(record, col, "Area"),
			ItemCode:      itemCode,
			ItemName:      fieldAt(record, col, "Item"),
			ElementCode:   elementCode,
			ElementName:   fieldAt(record, col, "Element"),
			Year:          strconv.Itoa(year),
			MonthsCode:    fieldAt(record, col, "Months Code"),
			Months:        fieldAt(record, col, "Months"),
			Unit:          fieldAt(record, col, "Unit"),
			Value:         fieldAt(record, col, "Value"),
			Flag:          fieldAt(record, col, "Flag"),
			CommoditySlug: ItemSlug[itemCode],
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no fao price rows after filter")
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
