package usda

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

const soapActionPerYear = `"http://www.fas.usda.gov/wsfaspsd/getDatabyCommodityPerYear"`

type psdRow struct {
	CommodityCode   string `json:"commodity_code"`
	CommodityName   string `json:"commodity_name"`
	CountryCode     string `json:"country_code"`
	CountryName     string `json:"country_name"`
	MarketingYear   string `json:"marketing_year"`
	CalendarYear    string `json:"calendar_year"`
	Month           string `json:"month"`
	AttributeID     string `json:"attribute_id"`
	AttributeName   string `json:"attribute_name"`
	UnitID          string `json:"unit_id"`
	UnitDescription string `json:"unit_description"`
	Value           string `json:"value"`
	CommoditySlug   string `json:"commodity_slug"`
}

// FetchPSDSnapshot downloads PSD rows for all countries and attributes by marketing year.
func (c *Client) FetchPSDSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	code := strings.TrimSpace(entry.PSDCommodityCode)
	if code == "" {
		return nil, "", fmt.Errorf("dataset %s missing psd_commodity_code", entry.DatasetID)
	}
	slug := strings.TrimSpace(entry.PSDCommoditySlug)
	if slug == "" {
		slug = psdSlugFromDataset(entry.DatasetID.String())
	}

	startYear, endYear, err := resolvePSDYearRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	merged := make(map[string]psdRow)
	yearsFetched := 0

	for year := startYear; year <= endYear; year++ {
		envelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <getDatabyCommodityPerYear xmlns="http://www.fas.usda.gov/wsfaspsd/">
      <strCommodityCode>%s</strCommodityCode>
      <strYear>%d</strYear>
    </getDatabyCommodityPerYear>
  </soap:Body>
</soap:Envelope>`, code, year)

		raw, err := c.PostSOAP(ctx, soapActionPerYear, envelope)
		if err != nil {
			continue
		}

		rows, err := parsePSDSoapResponse(raw)
		if err != nil {
			continue
		}
		for _, row := range rows {
			row.CommoditySlug = slug
			key := psdRowKey(row)
			merged[key] = row
		}
		yearsFetched++
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("usda psd returned no rows for %s (commodity=%s)", entry.DatasetID, code)
	}

	keys := make([]string, 0, len(merged))
	for key := range merged {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	out := make([]psdRow, 0, len(keys))
	for _, key := range keys {
		out = append(out, merged[key])
	}

	payload, err := json.Marshal(out)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s#getDatabyCommodityPerYear?commodity=%s&years=%d-%d&fetched=%d",
		psdSOAPEndpoint, code, startYear, endYear, yearsFetched)
	return payload, sourceURL, nil
}

func psdRowKey(row psdRow) string {
	return strings.Join([]string{
		row.CountryCode,
		row.MarketingYear,
		row.AttributeID,
		row.CalendarYear,
		row.Month,
	}, "|")
}

func psdSlugFromDataset(datasetID string) string {
	parts := strings.Split(datasetID, ".")
	if len(parts) != 2 {
		return ""
	}
	return strings.TrimPrefix(parts[1], "psd-")
}

func resolvePSDYearRange(entry catalog.RegistryEntry, fromDate string) (int, int, error) {
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

func parseYearHint(raw string) (int, bool) {
	raw = strings.TrimSpace(raw)
	if len(raw) >= 4 {
		if year, err := strconv.Atoi(raw[:4]); err == nil && year >= 1980 && year <= 2100 {
			return year, true
		}
	}
	return 0, false
}

// FlattenPSD converts merged PSD JSON into canonical bronze columns.
func FlattenPSD(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []psdRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse usda psd json: %w", err)
	}

	slug := strings.TrimSpace(entry.PSDCommoditySlug)
	if slug == "" {
		slug = psdSlugFromDataset(entry.DatasetID.String())
	}

	headers := []string{
		"commodity_code", "commodity_name", "commodity_slug",
		"country_code", "country_name",
		"marketing_year", "calendar_year", "month",
		"attribute_id", "attribute_name",
		"unit_id", "unit_description", "value",
	}
	out := make([][]string, 0, len(rows))

	for _, row := range rows {
		if strings.TrimSpace(row.CountryCode) == "" || strings.TrimSpace(row.MarketingYear) == "" {
			continue
		}
		rowSlug := strings.TrimSpace(row.CommoditySlug)
		if rowSlug == "" {
			rowSlug = slug
		}
		out = append(out, []string{
			row.CommodityCode,
			row.CommodityName,
			rowSlug,
			row.CountryCode,
			row.CountryName,
			row.MarketingYear,
			row.CalendarYear,
			row.Month,
			row.AttributeID,
			row.AttributeName,
			row.UnitID,
			row.UnitDescription,
			row.Value,
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no usda psd rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
