package fao

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultFFPICSVURL = "https://www.fao.org/fileadmin/templates/worldfood/Reports_and_docs/Food_price_indices_data.csv"

var ffpiIndexSlugs = map[string]string{
	"Food Price Index":  "food",
	"Meat Price Index":  "meat",
	"Dairy Price Index": "dairy",
	"Cereals Price Index": "cereals",
	"Oils Price Index":  "oils",
	"Sugar Price Index": "sugar",
}

type ffpiRow struct {
	RefMonth  string `json:"refmonth"`
	IndexSlug string `json:"index_slug"`
	IndexName string `json:"index_name"`
	Value     string `json:"value"`
	BasePeriod string `json:"base_period"`
}

// FetchFFPISnapshot downloads and parses the FAO monthly Food Price Index CSV.
func (c *Client) FetchFFPISnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = defaultFFPICSVURL
	}

	result, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	rows, err := parseFFPICSV(result.Body, entry)
	if err != nil {
		return nil, "", err
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("fao ffpi returned no rows for %s", entry.DatasetID)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].RefMonth != rows[j].RefMonth {
			return rows[i].RefMonth < rows[j].RefMonth
		}
		return rows[i].IndexSlug < rows[j].IndexSlug
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	return payload, sourceURL, nil
}

func parseFFPICSV(raw []byte, entry catalog.RegistryEntry) ([]ffpiRow, error) {
	reader := csv.NewReader(bytes.NewReader(raw))
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	minMonth := resolveFFPIMinMonth(entry)
	var header []string
	var rows []ffpiRow

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read fao ffpi row: %w", err)
		}
		if len(record) == 0 {
			continue
		}
		first := strings.TrimSpace(record[0])
		if first == "" || strings.HasPrefix(strings.ToUpper(first), "MONTHLY") {
			continue
		}
		if strings.EqualFold(first, "Date") {
			header = record
			continue
		}
		if len(header) == 0 {
			continue
		}

		refMonth, ok := parseFFPIMonth(first)
		if !ok || refMonth < minMonth {
			continue
		}

		for col := 1; col < len(header) && col < len(record); col++ {
			name := strings.TrimSpace(header[col])
			value := strings.TrimSpace(record[col])
			if name == "" || value == "" {
				continue
			}
			rows = append(rows, ffpiRow{
				RefMonth:   refMonth,
				IndexSlug:  ffpiSlug(name),
				IndexName:  name,
				Value:      value,
				BasePeriod: "2002-2004=100",
			})
		}
	}
	return rows, nil
}

func ffpiSlug(name string) string {
	if slug, ok := ffpiIndexSlugs[name]; ok {
		return slug
	}
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "_"))
}

func parseFFPIMonth(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	parts := strings.Split(raw, "-")
	if len(parts) != 2 {
		return "", false
	}
	monthMap := map[string]string{
		"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04",
		"May": "05", "Jun": "06", "Jul": "07", "Aug": "08",
		"Sep": "09", "Oct": "10", "Nov": "11", "Dec": "12",
	}
	mm, ok := monthMap[parts[0]]
	if !ok {
		return "", false
	}
	yy := parts[1]
	if len(yy) == 2 {
		if yy >= "90" {
			yy = "19" + yy
		} else {
			yy = "20" + yy
		}
	}
	return yy + "-" + mm + "-01", true
}

func resolveFFPIMinMonth(entry catalog.RegistryEntry) string {
	startYear := entry.PeriodStart
	if startYear == 0 {
		startYear = 2010
	}
	return fmt.Sprintf("%04d-01-01", startYear)
}

// FlattenFFPI converts merged FAO FFPI JSON into canonical bronze columns.
func FlattenFFPI(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []ffpiRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse fao ffpi json: %w", err)
	}

	headers := []string{"refmonth", "index_slug", "index_name", "value", "base_period"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RefMonth) == "" || strings.TrimSpace(row.Value) == "" {
			continue
		}
		out = append(out, []string{
			row.RefMonth,
			row.IndexSlug,
			row.IndexName,
			row.Value,
			row.BasePeriod,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no fao ffpi rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
