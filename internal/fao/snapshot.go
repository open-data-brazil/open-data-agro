package fao

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const (
	defaultProductionBulkURL = "https://bulks-faostat.fao.org/production/Production_Crops_Livestock_E_All_Data_(Normalized).zip"
	defaultProductionCSV       = "Production_Crops_Livestock_E_All_Data_(Normalized).csv"
	defaultTradeBulkURL        = "https://bulks-faostat.fao.org/production/Trade_Crops_Livestock_E_All_Data_(Normalized).zip"
	defaultTradeCSV            = "Trade_Crops_Livestock_E_All_Data_(Normalized).csv"
)

// FetchSnapshot downloads and filters FAOSTAT bulk data for a catalog entry.
func (c *Client) FetchSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	switch entry.DatasetID.String() {
	case "fao.producao-agro":
		return c.fetchAnnualBulkSnapshot(ctx, entry, fromDate, defaultProductionBulkURL, defaultProductionCSV, false)
	case "fao.comercio-agro":
		return c.fetchAnnualBulkSnapshot(ctx, entry, fromDate, defaultTradeBulkURL, defaultTradeCSV, false)
	case "fao.food-price-index":
		return c.FetchFFPISnapshot(ctx, entry)
	default:
		return c.FetchPricesSnapshot(ctx, entry, fromDate)
	}
}

func (c *Client) fetchAnnualBulkSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate, defaultBulkURL, defaultCSV string, withMonths bool) ([]byte, string, error) {
	bulkURL := strings.TrimSpace(entry.FAOBulkURL)
	if bulkURL == "" {
		bulkURL = defaultBulkURL
	}
	csvName := strings.TrimSpace(entry.FAOBulkCSV)
	if csvName == "" {
		csvName = defaultCSV
	}
	itemCodes := entry.FAOItemCodes
	if len(itemCodes) == 0 {
		itemCodes = defaultItemCodes
	}
	elementCodes := entry.FAOElementCodes
	if len(elementCodes) == 0 {
		return nil, "", fmt.Errorf("dataset %s missing fao_element_codes", entry.DatasetID)
	}

	startYear, endYear, err := resolveYearRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	result, err := c.Download(ctx, bulkURL)
	if err != nil {
		return nil, "", err
	}

	rows, err := parseFAOBulkZip(result.Body, csvName, codeSet(itemCodes), codeSet(elementCodes), startYear, endYear, withMonths)
	if err != nil {
		return nil, "", err
	}

	sort.Slice(rows, func(i, j int) bool {
		return bulkRowKey(rows[i]) < bulkRowKey(rows[j])
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s (items=%d, elements=%d, years=%d-%d, rows=%d)",
		bulkURL, len(itemCodes), len(elementCodes), startYear, endYear, len(rows))
	return payload, sourceURL, nil
}

// Flatten converts merged FAOSTAT JSON into canonical bronze columns.
func Flatten(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	switch entry.DatasetID.String() {
	case "fao.producao-agro", "fao.comercio-agro":
		return flattenAnnualBulk(entry, raw)
	case "fao.food-price-index":
		return FlattenFFPI(entry, raw)
	default:
		return FlattenPrices(entry, raw)
	}
}

func flattenAnnualBulk(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []bulkRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse fao bulk json: %w", err)
	}

	headers := []string{
		"area_code", "area_name", "item_code", "item_name", "commodity_slug",
		"element_code", "element_name", "year", "unit", "value", "flag",
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
			row.Unit,
			row.Value,
			row.Flag,
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no fao bulk rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
