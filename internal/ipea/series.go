package ipea

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

type odataSeriesResponse struct {
	Value []seriesValue `json:"value"`
}

type seriesValue struct {
	SeriesCode     string  `json:"SERCODIGO"`
	RefDate        string  `json:"VALDATA"`
	Value          float64 `json:"VALVALOR"`
	RegionLevel    string  `json:"NIVNOME"`
	TerritoryCode  string  `json:"TERCODIGO"`
}

type mergedSeriesPayload struct {
	Series []seriesValue `json:"series"`
}

// FetchSeriesSnapshot downloads and merges configured IPEA OData4 regional series.
func (c *Client) FetchSeriesSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	codes := entry.IpeaSeriesCodes
	if len(codes) == 0 {
		return nil, "", fmt.Errorf("dataset %s missing ipea_series_codes", entry.DatasetID)
	}

	var merged []seriesValue

	for _, code := range codes {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		requestURL := buildSeriesURL(code)

		result, err := c.Download(ctx, requestURL)
		if err != nil {
			return nil, "", fmt.Errorf("ipea fetch %s: %w", code, err)
		}

		var payload odataSeriesResponse
		if err := json.Unmarshal(result.Body, &payload); err != nil {
			return nil, "", fmt.Errorf("parse ipea response %s: %w", requestURL, err)
		}
		merged = append(merged, payload.Value...)
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("ipea returned no data rows for %s", entry.DatasetID)
	}

	sort.Slice(merged, func(i, j int) bool {
		return seriesRowKey(merged[i]) < seriesRowKey(merged[j])
	})

	out, err := json.Marshal(mergedSeriesPayload{Series: merged})
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s (series: %d, rows: %d)", odata4Base, len(codes), len(merged))
	return out, sourceURL, nil
}

func seriesRowKey(row seriesValue) string {
	ref := formatRefDate(row.RefDate)
	return strings.Join([]string{row.SeriesCode, ref, row.TerritoryCode}, "|")
}

func formatRefDate(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if len(raw) >= 10 {
		return raw[:10]
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t.Format("2006-01-02")
	}
	return raw
}

// FlattenSeries converts merged IPEA JSON into canonical bronze columns.
func FlattenSeries(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	rows, err := parseSeriesPayload(raw)
	if err != nil {
		return nil, nil, err
	}

	headers := []string{
		"series_code",
		"refdate",
		"value",
		"region_level",
		"territory_code",
	}

	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			strings.TrimSpace(row.SeriesCode),
			formatRefDate(row.RefDate),
			formatFloat(row.Value),
			strings.TrimSpace(row.RegionLevel),
			strings.TrimSpace(row.TerritoryCode),
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no ipea rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

func parseSeriesPayload(raw []byte) ([]seriesValue, error) {
	var wrapped mergedSeriesPayload
	if err := json.Unmarshal(raw, &wrapped); err == nil && len(wrapped.Series) > 0 {
		return wrapped.Series, nil
	}

	var direct odataSeriesResponse
	if err := json.Unmarshal(raw, &direct); err != nil {
		return nil, fmt.Errorf("parse ipea json: %w", err)
	}
	return direct.Value, nil
}

func formatFloat(value float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.4f", value), "0"), ".")
}
