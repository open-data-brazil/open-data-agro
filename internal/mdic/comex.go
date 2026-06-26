package mdic

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const generalEndpoint = "/general"

type comexRequest struct {
	Flow        string         `json:"flow"`
	MonthDetail bool           `json:"monthDetail"`
	Period      comexPeriod    `json:"period"`
	Filters     []comexFilter  `json:"filters"`
	Details     []string       `json:"details"`
	Metrics     []string       `json:"metrics"`
}

type comexPeriod struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type comexFilter struct {
	Filter string   `json:"filter"`
	Values []string `json:"values"`
}

type comexResponse struct {
	Success bool `json:"success"`
	Data    struct {
		List []comexRow `json:"list"`
	} `json:"data"`
}

type comexRow struct {
	CoNCM         string `json:"coNcm"`
	NCM           string `json:"ncm"`
	Year          string `json:"year"`
	MonthNumber   string `json:"monthNumber"`
	MetricFOB     string `json:"metricFOB"`
	MetricKG      string `json:"metricKG"`
}

type mergedRow struct {
	comexRow
	periodFrom string
	periodTo   string
}

// ProdutoSlug maps official NCM codes to canonical crop/product slugs.
var ProdutoSlug = map[string]string{
	"12019000": "soja",
	"10059000": "milho",
	"10019900": "trigo",
	"02013000": "carne_bovina",
	"02023000": "carne_bovina",
}

// FetchComexSnapshot downloads monthly export rows for configured NCM codes.
func (c *Client) FetchComexSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	flow := strings.TrimSpace(entry.ComexFlow)
	if flow == "" {
		flow = "export"
	}
	ncms := entry.ComexNCMs
	if len(ncms) == 0 {
		return nil, "", fmt.Errorf("dataset %s missing comex_ncms", entry.DatasetID)
	}

	start, end, err := resolveComexRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	merged := make(map[string]mergedRow)
	var chunkCount int

	for year := start.Year(); year <= end.Year(); year++ {
		periodFrom := fmt.Sprintf("%04d-01", year)
		periodTo := fmt.Sprintf("%04d-12", year)
		if year == start.Year() {
			periodFrom = fmt.Sprintf("%04d-%02d", year, int(start.Month()))
		}
		if year == end.Year() {
			periodTo = fmt.Sprintf("%04d-%02d", year, int(end.Month()))
		}

		reqBody, err := json.Marshal(comexRequest{
			Flow:        flow,
			MonthDetail: true,
			Period:      comexPeriod{From: periodFrom, To: periodTo},
			Filters:     []comexFilter{{Filter: "ncm", Values: ncms}},
			Details:     []string{"ncm"},
			Metrics:     []string{"metricFOB", "metricKG"},
		})
		if err != nil {
			return nil, "", err
		}

		raw, err := c.PostJSON(ctx, generalEndpoint, reqBody)
		if err != nil {
			return nil, "", fmt.Errorf("comex fetch %s-%s: %w", periodFrom, periodTo, err)
		}

		var resp comexResponse
		if err := json.Unmarshal(raw, &resp); err != nil {
			return nil, "", fmt.Errorf("parse comex response %s-%s: %w", periodFrom, periodTo, err)
		}
		if resp.Success == false && len(resp.Data.List) == 0 {
			continue
		}

		for _, row := range resp.Data.List {
			key := row.CoNCM + "|" + row.Year + "|" + row.MonthNumber
			merged[key] = mergedRow{comexRow: row, periodFrom: periodFrom, periodTo: periodTo}
		}
		chunkCount++

		select {
		case <-ctx.Done():
			return nil, "", ctx.Err()
		case <-time.After(400 * time.Millisecond):
		}
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("comex returned no data rows for %s", entry.DatasetID)
	}

	keys := make([]string, 0, len(merged))
	for key := range merged {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	out := make([]comexRow, 0, len(keys))
	for _, key := range keys {
		out = append(out, merged[key].comexRow)
	}

	payload, err := json.Marshal(out)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s%s (ncm=%d, chunks=%d)", comexAPIBase, generalEndpoint, len(ncms), chunkCount)
	return payload, sourceURL, nil
}

func resolveComexRange(entry catalog.RegistryEntry, fromDate string) (time.Time, time.Time, error) {
	startYear := entry.PeriodStart
	if startYear == 0 {
		startYear = 2015
	}
	start := time.Date(startYear, 1, 1, 0, 0, 0, 0, time.UTC)

	if raw := strings.TrimSpace(fromDate); raw != "" {
		parsed, err := parseISODate(raw)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid from date for %s: %w", entry.DatasetID, err)
		}
		if parsed.After(start) {
			start = parsed
		}
	}

	end := time.Now().UTC()
	if entry.PeriodEnd > 0 {
		candidate := time.Date(entry.PeriodEnd, 12, 1, 0, 0, 0, 0, time.UTC)
		if candidate.Before(end) {
			end = candidate
		}
	}
	if start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date range for %s", entry.DatasetID)
	}
	return start, end, nil
}

func parseISODate(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	layouts := []string{"2006-01-02", "2006-01"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unparseable date %q", raw)
}

// FlattenComex converts merged Comex JSON rows into canonical bronze columns.
func FlattenComex(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []comexRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse comex json: %w", err)
	}

	headers := []string{
		"co_ncm", "ncm_descricao", "produto_slug", "data",
		"valor_fob_usd", "quantidade_kg", "ano",
	}
	out := make([][]string, 0, len(rows))

	for _, row := range rows {
		coNCM := strings.TrimSpace(row.CoNCM)
		if coNCM == "" {
			continue
		}
		isoDate, err := monthToDate(row.Year, row.MonthNumber)
		if err != nil {
			continue
		}
		fob, err := normalizeMetric(row.MetricFOB)
		if err != nil {
			continue
		}
		kg, err := normalizeMetric(row.MetricKG)
		if err != nil {
			continue
		}
		slug := ProdutoSlug[coNCM]
		if slug == "" {
			slug = "outros"
		}
		out = append(out, []string{
			coNCM,
			strings.TrimSpace(row.NCM),
			slug,
			isoDate,
			fob,
			kg,
			isoDate[:4],
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no comex rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

func monthToDate(yearRaw, monthRaw string) (string, error) {
	year, err := strconv.Atoi(strings.TrimSpace(yearRaw))
	if err != nil || year < 1900 {
		return "", fmt.Errorf("invalid year %q", yearRaw)
	}
	month, err := strconv.Atoi(strings.TrimSpace(monthRaw))
	if err != nil || month < 1 || month > 12 {
		return "", fmt.Errorf("invalid month %q", monthRaw)
	}
	return fmt.Sprintf("%04d-%02d-01", year, month), nil
}

func normalizeMetric(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "0", nil
	}
	if _, err := strconv.ParseFloat(raw, 64); err != nil {
		return "", fmt.Errorf("invalid metric %q", raw)
	}
	return raw, nil
}
