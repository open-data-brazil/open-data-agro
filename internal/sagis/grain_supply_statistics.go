package sagis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

type grainRow struct {
	CommoditySlug   string `json:"commodity_slug"`
	MarketingYear   string `json:"marketing_year"`
	SupplyT         string `json:"supply_t"`
	DemandT         string `json:"demand_t"`
	OpeningStocksT  string `json:"opening_stocks_t"`
	ClosingStocksT  string `json:"closing_stocks_t"`
}

// Client downloads SAGIS grain supply statistics.
type Client struct{}

// NewClient creates a SAGIS client.
func NewClient() *Client { return &Client{} }

// FetchGrainSupplySnapshot returns SAGIS grain supply/demand rows (fixture until bulk URL stable).
func (c *Client) FetchGrainSupplySnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	_ = ctx
	rows, err := embeddedGrainSample()
	if err != nil {
		return nil, "", err
	}
	if path := strings.TrimSpace(os.Getenv("SAGIS_BULK_PATH")); path != "" {
		if parsed, parseErr := parseGrainJSONFile(path); parseErr == nil && len(parsed) > 0 {
			rows = parsed
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].CommoditySlug < rows[j].CommoditySlug
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = "https://www.sagis.org.za/"
	}
	return payload, sourceURL+" (fixture — bulk xlsx URL unstable 2026-06-26)", nil
}

func embeddedGrainSample() ([]grainRow, error) {
	return []grainRow{
		{CommoditySlug: "milho", MarketingYear: "2024", SupplyT: "14500000", DemandT: "11200000", OpeningStocksT: "4200000", ClosingStocksT: "7500000"},
		{CommoditySlug: "trigo", MarketingYear: "2024", SupplyT: "3200000", DemandT: "3400000", OpeningStocksT: "900000", ClosingStocksT: "700000"},
	}, nil
}

func parseGrainJSONFile(path string) ([]grainRow, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rows []grainRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// FlattenGrainSupply converts merged SAGIS JSON into canonical bronze columns.
func FlattenGrainSupply(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []grainRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse sagis json: %w", err)
	}

	headers := []string{"commodity_slug", "marketing_year", "supply_t", "demand_t", "opening_stocks_t", "closing_stocks_t"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.CommoditySlug) == "" {
			continue
		}
		out = append(out, []string{
			row.CommoditySlug, row.MarketingYear, row.SupplyT, row.DemandT,
			row.OpeningStocksT, row.ClosingStocksT,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no sagis rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
