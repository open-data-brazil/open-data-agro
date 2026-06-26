package igc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	xlsb "github.com/TsubasaBE/go-xlsb"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const goiBasePeriod = "2000-01=100"

// IndexSlug maps GOI sheet column headers to canonical slugs.
var IndexSlug = map[string]string{
	"IGC GOI":   "goi",
	"Wheat":     "wheat",
	"Maize":     "maize",
	"Soyabeans": "soybeans",
	"Rice":      "rice",
	"Barley":    "barley",
}

type goiRow struct {
	RefDate    string `json:"refdate"`
	IndexSlug  string `json:"index_slug"`
	IndexName  string `json:"index_name"`
	Value      string `json:"value"`
	BasePeriod string `json:"base_period"`
	Frequency  string `json:"frequency"`
}

// FetchGOISnapshot downloads the public IGC GOI xlsb workbook and extracts daily indices.
func (c *Client) FetchGOISnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	sourceURL, err := ResolveURL(entry)
	if err != nil {
		return nil, "", err
	}

	raw, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, "", fmt.Errorf("fetch igc goi: %w", err)
	}

	rows, err := parseGOIXLSB(raw, entry)
	if err != nil {
		return nil, "", err
	}

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	return payload, sourceURL, nil
}

func parseGOIXLSB(raw []byte, entry catalog.RegistryEntry) ([]goiRow, error) {
	wb, err := xlsb.OpenReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return nil, fmt.Errorf("open igc goi xlsb: %w", err)
	}
	defer func() { _ = wb.Close() }()

	sheetName := strings.TrimSpace(entry.XLSXSheet)
	if sheetName == "" {
		sheetName = "GOI & Indices"
	}

	ws, err := wb.SheetByName(sheetName)
	if err != nil {
		return nil, fmt.Errorf("igc sheet %q: %w", sheetName, err)
	}

	minDate := resolveGOIMinDate(entry)
	var header map[int]string
	var merged []goiRow

	for cells := range ws.Rows(false) {
		if len(cells) == 0 {
			continue
		}

		row := make(map[int]interface{})
		maxCol := 0
		for _, cell := range cells {
			row[cell.C] = cell.V
			if cell.C > maxCol {
				maxCol = cell.C
			}
		}

		if header == nil {
			if label, ok := row[0].(string); ok && strings.EqualFold(strings.TrimSpace(label), "DATE") {
				header = make(map[int]string)
				for col := 0; col <= maxCol; col++ {
					if name, ok := row[col].(string); ok {
						header[col] = strings.TrimSpace(name)
					}
				}
			}
			continue
		}

		refDate, ok := excelSerialToDate(row[0])
		if !ok || refDate < minDate {
			continue
		}

		for col, name := range header {
			if col == 0 || name == "" {
				continue
			}
			value, ok := numericValue(row[col])
			if !ok {
				continue
			}
			merged = append(merged, goiRow{
				RefDate:    refDate,
				IndexSlug:  indexSlug(name),
				IndexName:  name,
				Value:      strconv.FormatFloat(value, 'f', -1, 64),
				BasePeriod: goiBasePeriod,
				Frequency:  "daily",
			})
		}
	}

	if len(merged) == 0 {
		return nil, fmt.Errorf("no igc goi rows parsed for %s", entry.DatasetID)
	}

	sort.Slice(merged, func(i, j int) bool {
		if merged[i].RefDate != merged[j].RefDate {
			return merged[i].RefDate < merged[j].RefDate
		}
		return merged[i].IndexSlug < merged[j].IndexSlug
	})
	return merged, nil
}

func excelSerialToDate(raw interface{}) (string, bool) {
	switch v := raw.(type) {
	case float64:
		t, err := xlsb.ConvertDate(v)
		if err != nil {
			return "", false
		}
		return t.UTC().Format("2006-01-02"), true
	case int:
		t, err := xlsb.ConvertDate(float64(v))
		if err != nil {
			return "", false
		}
		return t.UTC().Format("2006-01-02"), true
	default:
		return "", false
	}
}

func numericValue(raw interface{}) (float64, bool) {
	switch v := raw.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

func indexSlug(name string) string {
	if slug, ok := IndexSlug[name]; ok {
		return slug
	}
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "_"))
}

func resolveGOIMinDate(entry catalog.RegistryEntry) string {
	startYear := entry.PeriodStart
	if startYear == 0 {
		startYear = 2010
	}
	return fmt.Sprintf("%04d-01-01", startYear)
}

// FlattenGOI converts merged GOI JSON into canonical bronze columns.
func FlattenGOI(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []goiRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse igc goi json: %w", err)
	}

	headers := []string{
		"refdate",
		"index_slug",
		"index_name",
		"value",
		"base_period",
		"frequency",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RefDate) == "" || strings.TrimSpace(row.Value) == "" {
			continue
		}
		out = append(out, []string{
			row.RefDate,
			row.IndexSlug,
			row.IndexName,
			row.Value,
			row.BasePeriod,
			row.Frequency,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no igc goi rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
