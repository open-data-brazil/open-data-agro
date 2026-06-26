package eia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const (
	eiaAPIBase      = "https://api.eia.gov/v2"
	defaultPageSize = 5000
)

// SeriesSlug maps EIA v1 series IDs to canonical commodity slugs.
var SeriesSlug = map[string]string{
	"PET.RWTC.D":  "wti_spot",
	"PET.RBRTE.D": "brent_spot",
}

// SeriesName maps EIA v1 series IDs to human-readable labels.
var SeriesName = map[string]string{
	"PET.RWTC.D":  "WTI Cushing OK Spot",
	"PET.RBRTE.D": "Europe Brent Spot",
}

type petroleumRow struct {
	SeriesID      string `json:"series_id"`
	SeriesName    string `json:"series_name"`
	CommoditySlug string `json:"commodity_slug"`
	RefDate       string `json:"refdate"`
	Unit          string `json:"unit"`
	Value         string `json:"value"`
	Frequency     string `json:"frequency"`
}

type eiaSeriesResponse struct {
	Response struct {
		Data []struct {
			Period string  `json:"period"`
			Value  float64 `json:"value"`
		} `json:"data"`
		Description  string `json:"description"`
		Frequency    string `json:"frequency"`
		Units        string `json:"units"`
		Total        string `json:"total"`
		DateFormat   string `json:"dateFormat"`
	} `json:"response"`
}

// FetchPetroleumSnapshot downloads daily petroleum spot prices via EIA API v2 seriesid route.
func (c *Client) FetchPetroleumSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	apiKey := strings.TrimSpace(os.Getenv("EIA_API_KEY"))
	if apiKey == "" {
		return nil, "", fmt.Errorf("EIA_API_KEY is required for %s", entry.DatasetID)
	}

	seriesIDs := entry.EIASeriesIDs
	if len(seriesIDs) == 0 {
		seriesIDs = defaultSeriesIDs()
	}

	start, end, err := resolvePetroleumDateRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	var merged []petroleumRow
	var requestURLs []string

	for _, seriesID := range seriesIDs {
		seriesID = strings.TrimSpace(seriesID)
		if seriesID == "" {
			continue
		}

		offset := 0
		for {
			requestURL := buildSeriesURL(seriesID, apiKey, start, end, offset)
			requestURLs = append(requestURLs, requestURL)

			result, err := c.Download(ctx, requestURL)
			if err != nil {
				return nil, "", fmt.Errorf("eia fetch %s: %w", seriesID, err)
			}

			rows, total, _, _, err := parseSeriesResponse(seriesID, result.Body)
			if err != nil {
				return nil, "", fmt.Errorf("parse eia response %s: %w", seriesID, err)
			}
			merged = append(merged, rows...)

			offset += len(rows)
			if len(rows) == 0 || offset >= total {
				break
			}
		}
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("eia returned no data rows for %s", entry.DatasetID)
	}

	sort.Slice(merged, func(i, j int) bool {
		if merged[i].SeriesID != merged[j].SeriesID {
			return merged[i].SeriesID < merged[j].SeriesID
		}
		return merged[i].RefDate < merged[j].RefDate
	})

	payload, err := json.Marshal(merged)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s/seriesid (series=%d, chunks=%d)", eiaAPIBase, len(seriesIDs), len(requestURLs))
	return payload, sourceURL, nil
}

func buildSeriesURL(seriesID, apiKey, start, end string, offset int) string {
	values := url.Values{}
	values.Set("api_key", apiKey)
	values.Set("start", start)
	values.Set("end", end)
	values.Set("length", strconv.Itoa(defaultPageSize))
	values.Set("offset", strconv.Itoa(offset))
	values.Set("sort[0][column]", "period")
	values.Set("sort[0][direction]", "asc")
	return fmt.Sprintf("%s/seriesid/%s?%s", eiaAPIBase, url.PathEscape(seriesID), values.Encode())
}

func parseSeriesResponse(seriesID string, raw []byte) ([]petroleumRow, int, string, string, error) {
	var payload eiaSeriesResponse
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, 0, "", "", err
	}

	total, _ := strconv.Atoi(strings.TrimSpace(payload.Response.Total))
	frequency := strings.TrimSpace(payload.Response.Frequency)
	unit := strings.TrimSpace(payload.Response.Units)
	if unit == "" {
		unit = "dollars per barrel"
	}

	name := seriesLabel(seriesID, payload.Response.Description)
	slug := seriesSlug(seriesID)

	rows := make([]petroleumRow, 0, len(payload.Response.Data))
	for _, item := range payload.Response.Data {
		period := strings.TrimSpace(item.Period)
		if period == "" {
			continue
		}
		rows = append(rows, petroleumRow{
			SeriesID:      seriesID,
			SeriesName:    name,
			CommoditySlug: slug,
			RefDate:       period,
			Unit:          unit,
			Value:         strconv.FormatFloat(item.Value, 'f', -1, 64),
			Frequency:     frequency,
		})
	}

	return rows, total, frequency, unit, nil
}

func seriesLabel(seriesID, description string) string {
	if name, ok := SeriesName[seriesID]; ok {
		return name
	}
	if strings.TrimSpace(description) != "" {
		return strings.TrimSpace(description)
	}
	return seriesID
}

func seriesSlug(seriesID string) string {
	if slug, ok := SeriesSlug[seriesID]; ok {
		return slug
	}
	return strings.ToLower(strings.ReplaceAll(seriesID, ".", "_"))
}

func defaultSeriesIDs() []string {
	return []string{"PET.RWTC.D", "PET.RBRTE.D"}
}

func resolvePetroleumDateRange(entry catalog.RegistryEntry, fromDate string) (string, string, error) {
	startYear := entry.PeriodStart
	if startYear == 0 {
		startYear = 2010
	}
	start := fmt.Sprintf("%d-01-01", startYear)

	fromDate = strings.TrimSpace(fromDate)
	if fromDate != "" {
		parsed, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			return "", "", fmt.Errorf("invalid fromDate %q: %w", fromDate, err)
		}
		start = parsed.Format("2006-01-02")
	}

	end := time.Now().UTC().Format("2006-01-02")
	return start, end, nil
}

// FlattenPetroleum converts merged EIA petroleum JSON into canonical bronze columns.
func FlattenPetroleum(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []petroleumRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse eia petroleum json: %w", err)
	}

	headers := []string{
		"series_id",
		"series_name",
		"commodity_slug",
		"refdate",
		"unit",
		"value",
		"frequency",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RefDate) == "" || strings.TrimSpace(row.Value) == "" {
			continue
		}
		out = append(out, []string{
			row.SeriesID,
			row.SeriesName,
			row.CommoditySlug,
			row.RefDate,
			row.Unit,
			row.Value,
			row.Frequency,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no petroleum rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

// ResolveURL validates the catalog base URL for an EIA dataset.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return "", fmt.Errorf("dataset %s has no source_url", entry.DatasetID)
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid source_url for %s: %w", entry.DatasetID, err)
	}
	if parsed.Scheme != "https" {
		return "", fmt.Errorf("source_url for %s must use https", entry.DatasetID)
	}
	if !strings.EqualFold(parsed.Host, "api.eia.gov") {
		return "", fmt.Errorf("source_url for %s must be on api.eia.gov", entry.DatasetID)
	}
	return parsed.String(), nil
}
