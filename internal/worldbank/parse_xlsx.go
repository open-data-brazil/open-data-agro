package worldbank

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	defaultPinkSheetSheet = "Monthly Prices"
	pinkSheetHeaderRow    = 4
	pinkSheetUnitRow      = 5
	pinkSheetDataStartRow = 6
)

func parsePinkSheetXLSX(raw []byte, sheetName string, seriesFilter map[string]struct{}, startMonth, endMonth string) ([]pinkSheetRow, error) {
	if strings.TrimSpace(sheetName) == "" {
		sheetName = defaultPinkSheetSheet
	}

	book, err := excelize.OpenReader(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("open pink sheet xlsx: %w", err)
	}
	defer func() { _ = book.Close() }()

	table, err := book.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("read sheet %q: %w", sheetName, err)
	}
	if len(table) <= pinkSheetDataStartRow {
		return nil, fmt.Errorf("pink sheet sheet %q too short", sheetName)
	}

	headers := table[pinkSheetHeaderRow]
	units := table[pinkSheetUnitRow]
	columns := make([]seriesColumn, 0, len(headers))
	for i := 1; i < len(headers); i++ {
		name := strings.TrimSpace(headers[i])
		if name == "" {
			continue
		}
		if len(seriesFilter) > 0 {
			if _, ok := seriesFilter[name]; !ok {
				continue
			}
		}
		unit := ""
		if i < len(units) {
			unit = strings.TrimSpace(units[i])
		}
		columns = append(columns, seriesColumn{Name: name, Index: i, Unit: unit})
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("no matching pink sheet series columns")
	}

	var rows []pinkSheetRow
	for _, record := range table[pinkSheetDataStartRow:] {
		if len(record) == 0 {
			continue
		}
		refmonth, ok := parseRefMonth(record[0])
		if !ok {
			continue
		}
		if refmonth < startMonth || refmonth > endMonth {
			continue
		}

		for _, col := range columns {
			if col.Index >= len(record) {
				continue
			}
			value := normalizeCell(record[col.Index])
			if value == "" {
				continue
			}
			rows = append(rows, pinkSheetRow{
				RefMonth:      refmonth,
				SeriesName:    col.Name,
				Unit:          col.Unit,
				Value:         value,
				CommoditySlug: SeriesSlug[col.Name],
			})
		}
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no pink sheet rows after filter")
	}
	return rows, nil
}

type seriesColumn struct {
	Name  string
	Index int
	Unit  string
}

func parseRefMonth(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if len(raw) != 7 || raw[4] != 'M' {
		return "", false
	}
	year, err := strconv.Atoi(raw[:4])
	if err != nil || year < 1960 {
		return "", false
	}
	month, err := strconv.Atoi(raw[5:])
	if err != nil || month < 1 || month > 12 {
		return "", false
	}
	return fmt.Sprintf("%04d-%02d", year, month), true
}

func normalizeCell(raw string) string {
	text := strings.TrimSpace(raw)
	if text == "" || text == "…" || strings.EqualFold(text, "na") {
		return ""
	}
	return text
}

func seriesNameSet(names []string) map[string]struct{} {
	out := make(map[string]struct{}, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		out[name] = struct{}{}
	}
	return out
}
