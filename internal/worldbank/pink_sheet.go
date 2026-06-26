package worldbank

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

const defaultPinkSheetURL = "https://thedocs.worldbank.org/en/doc/74e8be41ceb20fa0da750cda2f6b9e4e-0050012026/related/CMO-Historical-Data-Monthly.xlsx"

// SeriesSlug maps Pink Sheet series names to canonical commodity slugs.
var SeriesSlug = map[string]string{
	"Soybeans":           "soja",
	"Maize":              "milho",
	"Wheat, US SRW":      "trigo",
	"Beef **":            "carne_bovina",
	"Crude oil, average": "petroleo",
}

var defaultSeriesNames = []string{
	"Soybeans",
	"Maize",
	"Wheat, US SRW",
	"Beef **",
	"Crude oil, average",
}

var defaultAgIndexSeries = []string{
	"Agriculture **",
	"Food **",
	"Grains",
	"Oils & Meals",
	"Beverages",
	"Raw Materials",
	"Timber",
	"Other Food **",
	"Other Raw Mat.",
	"Fertilizers **",
}

// AgIndexSlug maps Pink Sheet agriculture index names to canonical slugs.
var AgIndexSlug = map[string]string{
	"Agriculture **":  "agriculture",
	"Food **":         "food",
	"Grains":          "grains",
	"Oils & Meals":    "oleos_e_farelo",
	"Beverages":       "beverages",
	"Raw Materials":   "raw_materials",
	"Timber":          "timber",
	"Other Food **":   "other_food",
	"Other Raw Mat.":  "other_raw_materials",
	"Fertilizers **":  "fertilizers",
}

type pinkSheetRow struct {
	RefMonth      string `json:"refmonth"`
	SeriesName    string `json:"series_name"`
	Unit          string `json:"unit"`
	Value         string `json:"value"`
	CommoditySlug string `json:"commodity_slug"`
}

// FetchSnapshot downloads and unpivots World Bank Pink Sheet data for a catalog entry.
func (c *Client) FetchSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	switch entry.DatasetID.String() {
	case "worldbank.ag-indices":
		return c.FetchAgIndicesSnapshot(ctx, entry, fromDate)
	default:
		return c.FetchPinkSheetSnapshot(ctx, entry, fromDate)
	}
}

// FetchPinkSheetSnapshot downloads and unpivots World Bank Pink Sheet monthly prices.
func (c *Client) FetchPinkSheetSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	sourceURL := strings.TrimSpace(entry.WorldBankPinkSheetURL)
	if sourceURL == "" {
		sourceURL = defaultPinkSheetURL
	}
	sheet := strings.TrimSpace(entry.WorldBankPinkSheetSheet)
	seriesNames := entry.WorldBankSeriesNames
	if len(seriesNames) == 0 {
		seriesNames = defaultSeriesNames
	}

	startMonth, endMonth, err := resolveMonthRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	result, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	rows, err := parsePinkSheetXLSX(result.Body, sheet, seriesNameSet(seriesNames), startMonth, endMonth)
	if err != nil {
		return nil, "", err
	}

	sort.Slice(rows, func(i, j int) bool {
		return pinkSheetRowKey(rows[i]) < pinkSheetRowKey(rows[j])
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}

	metaURL := fmt.Sprintf("%s (sheet=%s, series=%d, months=%s-%s, rows=%d)",
		sourceURL, sheetNameOrDefault(sheet), len(seriesNames), startMonth, endMonth, len(rows))
	return payload, metaURL, nil
}

// FetchAgIndicesSnapshot downloads agriculture sub-indices from the Pink Sheet workbook.
func (c *Client) FetchAgIndicesSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	sourceURL := strings.TrimSpace(entry.WorldBankPinkSheetURL)
	if sourceURL == "" {
		sourceURL = defaultPinkSheetURL
	}
	sheet := strings.TrimSpace(entry.WorldBankPinkSheetSheet)
	if sheet == "" {
		sheet = defaultAgIndicesSheet
	}
	seriesNames := entry.WorldBankSeriesNames
	if len(seriesNames) == 0 {
		seriesNames = defaultAgIndexSeries
	}

	startMonth, endMonth, err := resolveMonthRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	result, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	rows, err := parseAgIndicesXLSX(result.Body, sheet, seriesNameSet(seriesNames), startMonth, endMonth)
	if err != nil {
		return nil, "", err
	}

	sort.Slice(rows, func(i, j int) bool {
		return pinkSheetRowKey(rows[i]) < pinkSheetRowKey(rows[j])
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}

	metaURL := fmt.Sprintf("%s (sheet=%s, series=%d, months=%s-%s, rows=%d)",
		sourceURL, sheet, len(seriesNames), startMonth, endMonth, len(rows))
	return payload, metaURL, nil
}

func sheetNameOrDefault(sheet string) string {
	if strings.TrimSpace(sheet) == "" {
		return defaultPinkSheetSheet
	}
	return sheet
}

func pinkSheetRowKey(row pinkSheetRow) string {
	return strings.Join([]string{row.RefMonth, row.SeriesName}, "|")
}

func resolveMonthRange(entry catalog.RegistryEntry, fromDate string) (string, string, error) {
	startYear := entry.PeriodStart
	if startYear == 0 {
		startYear = 2010
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

func parseMonthHint(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if len(raw) >= 7 && raw[4] == '-' {
		if _, err := time.Parse("2006-01", raw[:7]); err == nil {
			return raw[:7], true
		}
	}
	if len(raw) >= 4 {
		if year, err := strconv.Atoi(raw[:4]); err == nil && year >= 1960 {
			return fmt.Sprintf("%04d-01", year), true
		}
	}
	return "", false
}

// Flatten converts merged World Bank JSON into canonical bronze columns.
func Flatten(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	return FlattenPinkSheet(entry, raw)
}

// FlattenPinkSheet converts merged Pink Sheet JSON into canonical bronze columns.
func FlattenPinkSheet(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []pinkSheetRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse worldbank pink sheet json: %w", err)
	}

	headers := []string{
		"refmonth", "series_name", "commodity_slug", "unit", "value",
	}
	out := make([][]string, 0, len(rows))

	for _, row := range rows {
		if strings.TrimSpace(row.RefMonth) == "" || strings.TrimSpace(row.SeriesName) == "" {
			continue
		}
		slug := strings.TrimSpace(row.CommoditySlug)
		if slug == "" {
			slug = SeriesSlug[row.SeriesName]
			if slug == "" {
				slug = AgIndexSlug[row.SeriesName]
			}
		}
		out = append(out, []string{
			row.RefMonth,
			row.SeriesName,
			slug,
			row.Unit,
			row.Value,
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no worldbank pink sheet rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
