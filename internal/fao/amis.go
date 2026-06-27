package fao

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

type amisRow struct {
	CommoditySlug  string `json:"commodity_slug"`
	RefMonth       string `json:"refmonth"`
	IndicatorSlug  string `json:"indicator_slug"`
	Value          string `json:"value"`
	Unit           string `json:"unit"`
}

// FetchAMISMarketMonitorSnapshot returns FAO AMIS market monitor indicators (fixture — SPA export).
func (c *Client) FetchAMISMarketMonitorSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	_ = ctx
	rows, err := embeddedAMISample()
	if err != nil {
		return nil, "", err
	}
	if path := strings.TrimSpace(os.Getenv("AMIS_BULK_PATH")); path != "" {
		if parsed, parseErr := parseAMISJSONFile(path); parseErr == nil && len(parsed) > 0 {
			rows = parsed
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].CommoditySlug+rows[i].RefMonth < rows[j].CommoditySlug+rows[j].RefMonth
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = "https://www.amis-outlook.org/#/market-database"
	}
	return payload, sourceURL + " (fixture — market database CSV behind SPA 2026-06-26)", nil
}

func embeddedAMISample() ([]amisRow, error) {
	return []amisRow{
		{CommoditySlug: "milho", RefMonth: "2024-01-01", IndicatorSlug: "price_index", Value: "118.5", Unit: "index_2007_2008=100"},
		{CommoditySlug: "soja", RefMonth: "2024-01-01", IndicatorSlug: "price_index", Value: "112.3", Unit: "index_2007_2008=100"},
	}, nil
}

func parseAMISJSONFile(path string) ([]amisRow, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rows []amisRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// FlattenAMISMarketMonitor converts merged AMIS JSON into canonical bronze columns.
func FlattenAMISMarketMonitor(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []amisRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse fao amis json: %w", err)
	}

	headers := []string{"commodity_slug", "refmonth", "indicator_slug", "value", "unit"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.CommoditySlug) == "" {
			continue
		}
		out = append(out, []string{
			row.CommoditySlug, row.RefMonth, row.IndicatorSlug, row.Value, row.Unit,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no fao amis rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
