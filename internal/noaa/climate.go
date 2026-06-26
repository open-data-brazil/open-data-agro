package noaa

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

const (
	defaultONIURL        = "https://www.cpc.ncep.noaa.gov/data/indices/oni.ascii.txt"
	defaultGlobalTempURL = "https://www.ncei.noaa.gov/access/monitoring/climate-at-a-glance/global/time-series/globe/land_ocean/0/0/%d-%d.csv"
)

type ensoRow struct {
	RefYear    string `json:"refyear"`
	SeasonCode string `json:"season_code"`
	SSTTotal   string `json:"sst_total"`
	Anomaly    string `json:"anomaly"`
	IndexName  string `json:"index_name"`
}

type globalTempRow struct {
	RefMonth   string `json:"refmonth"`
	Anomaly    string `json:"anomaly"`
	Unit       string `json:"unit"`
	BasePeriod string `json:"base_period"`
	IndexName  string `json:"index_name"`
}

// FetchENSOSnapshot downloads NOAA ONI (Oceanic Niño Index) ASCII data.
func (c *Client) FetchENSOSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	sourceURL := strings.TrimSpace(entry.NOAAIndexURL)
	if sourceURL == "" {
		sourceURL = defaultONIURL
	}

	startYear, endYear, err := resolveYearRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	result, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	rows, err := parseONIASCII(result.Body, startYear, endYear)
	if err != nil {
		return nil, "", err
	}

	sort.Slice(rows, func(i, j int) bool {
		return ensoRowKey(rows[i]) < ensoRowKey(rows[j])
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	meta := fmt.Sprintf("%s (years=%d-%d, rows=%d)", sourceURL, startYear, endYear, len(rows))
	return payload, meta, nil
}

// FetchGlobalTempSnapshot downloads NCEI global land+ocean temperature anomaly CSV.
func (c *Client) FetchGlobalTempSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	startYear, endYear, err := resolveYearRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	sourceURL := strings.TrimSpace(entry.NOAAIndexURL)
	if sourceURL == "" {
		sourceURL = fmt.Sprintf(defaultGlobalTempURL, startYear, endYear)
	}

	result, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	startMonth, endMonth, err := resolveMonthRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	rows, err := parseGlobalTempCSV(result.Body, startMonth, endMonth)
	if err != nil {
		return nil, "", err
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].RefMonth < rows[j].RefMonth
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	meta := fmt.Sprintf("%s (months=%s-%s, rows=%d)", sourceURL, startMonth, endMonth, len(rows))
	return payload, meta, nil
}

func ensoRowKey(row ensoRow) string {
	return row.RefYear + "|" + row.SeasonCode
}

func resolveYearRange(entry catalog.RegistryEntry, fromDate string) (int, int, error) {
	startYear := entry.PeriodStart
	if startYear == 0 {
		startYear = 2010
	}
	if raw := strings.TrimSpace(fromDate); raw != "" {
		if year, ok := parseYearHint(raw); ok && year > startYear {
			startYear = year
		}
	}
	endYear := time.Now().UTC().Year()
	if entry.PeriodEnd > 0 && entry.PeriodEnd < endYear {
		endYear = entry.PeriodEnd
	}
	if startYear > endYear {
		return 0, 0, fmt.Errorf("invalid year range for %s", entry.DatasetID)
	}
	return startYear, endYear, nil
}

func resolveMonthRange(entry catalog.RegistryEntry, fromDate string) (string, string, error) {
	startYear, _, err := resolveYearRange(entry, fromDate)
	if err != nil {
		return "", "", err
	}
	startMonth := fmt.Sprintf("%04d-01", startYear)
	if raw := strings.TrimSpace(fromDate); raw != "" {
		if month, ok := parseMonthHint(raw); ok && month > startMonth {
			startMonth = month
		}
	}
	end := time.Now().UTC()
	endMonth := fmt.Sprintf("%04d-%02d", end.Year(), int(end.Month()))
	if entry.PeriodEnd > 0 {
		candidate := fmt.Sprintf("%04d-12", entry.PeriodEnd)
		if candidate < endMonth {
			endMonth = candidate
		}
	}
	if startMonth > endMonth {
		return "", "", fmt.Errorf("invalid month range for %s", entry.DatasetID)
	}
	return startMonth, endMonth, nil
}

func parseYearHint(raw string) (int, bool) {
	raw = strings.TrimSpace(raw)
	if len(raw) >= 4 {
		if year, err := strconv.Atoi(raw[:4]); err == nil && year >= 1950 {
			return year, true
		}
	}
	return 0, false
}

func parseMonthHint(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if len(raw) >= 7 && raw[4] == '-' {
		if _, err := time.Parse("2006-01", raw[:7]); err == nil {
			return raw[:7], true
		}
	}
	if year, ok := parseYearHint(raw); ok {
		return fmt.Sprintf("%04d-01", year), true
	}
	return "", false
}

// FlattenENSO converts merged ONI JSON into canonical bronze columns.
func FlattenENSO(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []ensoRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse noaa enso json: %w", err)
	}
	headers := []string{"refyear", "season_code", "sst_total", "anomaly", "index_name"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RefYear) == "" || strings.TrimSpace(row.SeasonCode) == "" {
			continue
		}
		out = append(out, []string{row.RefYear, row.SeasonCode, row.SSTTotal, row.Anomaly, row.IndexName})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no noaa enso rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

// FlattenGlobalTemp converts merged global temperature JSON into bronze columns.
func FlattenGlobalTemp(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []globalTempRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse noaa global temp json: %w", err)
	}
	headers := []string{"refmonth", "anomaly", "unit", "base_period", "index_name"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RefMonth) == "" {
			continue
		}
		out = append(out, []string{row.RefMonth, row.Anomaly, row.Unit, row.BasePeriod, row.IndexName})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no noaa global temp rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

// FetchClimateSnapshot routes to the correct NOAA index fetcher by dataset ID.
func (c *Client) FetchClimateSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	switch entry.DatasetID.String() {
	case "noaa.enso-indices":
		return c.FetchENSOSnapshot(ctx, entry, fromDate)
	case "noaa.global-temp-anomaly":
		return c.FetchGlobalTempSnapshot(ctx, entry, fromDate)
	default:
		return nil, "", fmt.Errorf("unsupported noaa dataset %s", entry.DatasetID)
	}
}

// FlattenClimate routes flattening by dataset ID.
func FlattenClimate(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	switch entry.DatasetID.String() {
	case "noaa.enso-indices":
		return FlattenENSO(entry, raw)
	case "noaa.global-temp-anomaly":
		return FlattenGlobalTemp(entry, raw)
	default:
		return nil, nil, fmt.Errorf("unsupported noaa dataset %s", entry.DatasetID)
	}
}
