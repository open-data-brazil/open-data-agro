package fao

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

const defaultPricesBulkURL = "https://bulks-faostat.fao.org/production/Prices_E_All_Data_(Normalized).zip"

// ItemSlug maps FAOSTAT item codes to canonical commodity slugs.
var ItemSlug = map[string]string{
	"236": "soja",
	"56":  "milho",
	"15":  "trigo",
	"867": "carne_bovina",
}

var defaultItemCodes = []string{"236", "56", "15", "867"}
var defaultElementCodes = []string{"5532", "5539"}

type priceRow struct {
	AreaCode      string `json:"area_code"`
	AreaName      string `json:"area_name"`
	ItemCode      string `json:"item_code"`
	ItemName      string `json:"item_name"`
	ElementCode   string `json:"element_code"`
	ElementName   string `json:"element_name"`
	Year          string `json:"year"`
	MonthsCode    string `json:"months_code"`
	Months        string `json:"months"`
	Unit          string `json:"unit"`
	Value         string `json:"value"`
	Flag          string `json:"flag"`
	CommoditySlug string `json:"commodity_slug"`
}

// FetchPricesSnapshot downloads and filters FAOSTAT producer prices bulk data.
func (c *Client) FetchPricesSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	bulkURL := strings.TrimSpace(entry.FAOBulkURL)
	if bulkURL == "" {
		bulkURL = defaultPricesBulkURL
	}
	csvName := strings.TrimSpace(entry.FAOBulkCSV)
	itemCodes := entry.FAOItemCodes
	if len(itemCodes) == 0 {
		itemCodes = defaultItemCodes
	}
	elementCodes := entry.FAOElementCodes
	if len(elementCodes) == 0 {
		elementCodes = defaultElementCodes
	}

	startYear, endYear, err := resolveYearRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	result, err := c.Download(ctx, bulkURL)
	if err != nil {
		return nil, "", err
	}

	rows, err := parsePricesBulkZip(result.Body, csvName, codeSet(itemCodes), codeSet(elementCodes), startYear, endYear)
	if err != nil {
		return nil, "", err
	}

	sort.Slice(rows, func(i, j int) bool {
		return priceRowKey(rows[i]) < priceRowKey(rows[j])
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s (items=%d, elements=%d, years=%d-%d, rows=%d)",
		bulkURL, len(itemCodes), len(elementCodes), startYear, endYear, len(rows))
	return payload, sourceURL, nil
}

func priceRowKey(row priceRow) string {
	return strings.Join([]string{
		row.AreaCode, row.ItemCode, row.ElementCode, row.Year, row.MonthsCode,
	}, "|")
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

func parseYearHint(raw string) (int, bool) {
	raw = strings.TrimSpace(raw)
	if len(raw) >= 4 {
		if year, err := strconv.Atoi(raw[:4]); err == nil && year >= 1960 && year <= 2100 {
			return year, true
		}
	}
	return 0, false
}

// FlattenPrices converts merged FAOSTAT prices JSON into canonical bronze columns.
func FlattenPrices(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []priceRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse fao prices json: %w", err)
	}

	headers := []string{
		"area_code", "area_name", "item_code", "item_name", "commodity_slug",
		"element_code", "element_name", "year", "months_code", "months",
		"unit", "value", "flag",
	}
	out := make([][]string, 0, len(rows))

	for _, row := range rows {
		if strings.TrimSpace(row.AreaCode) == "" || strings.TrimSpace(row.Year) == "" {
			continue
		}
		slug := strings.TrimSpace(row.CommoditySlug)
		if slug == "" {
			slug = ItemSlug[row.ItemCode]
		}
		out = append(out, []string{
			row.AreaCode,
			row.AreaName,
			row.ItemCode,
			row.ItemName,
			slug,
			row.ElementCode,
			row.ElementName,
			row.Year,
			row.MonthsCode,
			row.Months,
			row.Unit,
			row.Value,
			row.Flag,
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no fao price rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
