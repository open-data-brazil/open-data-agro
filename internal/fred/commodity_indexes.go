package fred

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

const defaultCommodityIndexURL = "https://fred.stlouisfed.org/graph/fredgraph.csv?id=PALLFNFINDEXM"

var defaultSeriesIDs = []string{"PALLFNFINDEXM", "PPIACO", "WPSFD49207"}

type indexRow struct {
	SeriesID string `json:"series_id"`
	RefMonth string `json:"refmonth"`
	Value    string `json:"value"`
}

// FetchCommodityIndexesSnapshot downloads FRED commodity index CSV series.
func (c *Client) FetchCommodityIndexesSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	seriesIDs := entry.WorldBankSeriesNames
	if len(seriesIDs) == 0 {
		seriesIDs = defaultSeriesIDs[:1]
	}

	var merged []indexRow
	var sourceURLs []string
	minMonth := resolveMinMonth(entry)

	for _, seriesID := range seriesIDs {
		sourceURL := buildSeriesURL(entry, seriesID)
		raw, err := c.download(ctx, sourceURL)
		if err != nil {
			return nil, "", fmt.Errorf("fetch fred %s: %w", seriesID, err)
		}
		rows, err := parseFREDCSV(raw, seriesID, minMonth)
		if err != nil {
			return nil, "", err
		}
		merged = append(merged, rows...)
		sourceURLs = append(sourceURLs, sourceURL)
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("fred commodity indexes returned no rows for %s", entry.DatasetID)
	}

	sort.Slice(merged, func(i, j int) bool {
		if merged[i].SeriesID != merged[j].SeriesID {
			return merged[i].SeriesID < merged[j].SeriesID
		}
		return merged[i].RefMonth < merged[j].RefMonth
	})

	payload, err := json.Marshal(merged)
	if err != nil {
		return nil, "", err
	}
	return payload, strings.Join(sourceURLs, "; "), nil
}

func buildSeriesURL(entry catalog.RegistryEntry, seriesID string) string {
	if u := strings.TrimSpace(entry.SourceURL); u != "" && strings.Contains(u, seriesID) {
		return u
	}
	return fmt.Sprintf("https://fred.stlouisfed.org/graph/fredgraph.csv?id=%s", seriesID)
}

func resolveMinMonth(entry catalog.RegistryEntry) string {
	startYear := entry.PeriodStart
	if startYear == 0 {
		startYear = 2010
	}
	return fmt.Sprintf("%04d-01-01", startYear)
}

func parseFREDCSV(raw []byte, seriesID, minMonth string) ([]indexRow, error) {
	reader := csv.NewReader(bytes.NewReader(raw))
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read fred header: %w", err)
	}
	valueCol := seriesID
	if len(header) == 2 {
		valueCol = strings.TrimSpace(header[1])
	}

	var rows []indexRow
	for {
		record, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, fmt.Errorf("read fred row: %w", readErr)
		}
		if len(record) < 2 {
			continue
		}
		refMonth := normalizeMonth(strings.TrimSpace(record[0]))
		value := strings.TrimSpace(record[1])
		if refMonth == "" || value == "" || refMonth < minMonth {
			continue
		}
		rows = append(rows, indexRow{
			SeriesID: valueCol,
			RefMonth: refMonth,
			Value:    value,
		})
	}
	return rows, nil
}

func normalizeMonth(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) == 7 {
		return raw + "-01"
	}
	return raw
}

// FlattenCommodityIndexes converts merged FRED JSON into canonical bronze columns.
func FlattenCommodityIndexes(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []indexRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse fred json: %w", err)
	}

	headers := []string{"series_id", "refmonth", "value"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RefMonth) == "" || strings.TrimSpace(row.Value) == "" {
			continue
		}
		out = append(out, []string{row.SeriesID, row.RefMonth, row.Value})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no fred rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
