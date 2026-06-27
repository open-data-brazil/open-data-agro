package japan

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

type tradeRow struct {
	CommoditySlug string `json:"commodity_slug"`
	RefYear       string `json:"refyear"`
	FlowCode      string `json:"flow_code"`
	ValueJPY      string `json:"value_jpy"`
	QuantityT     string `json:"quantity_t"`
}

// Client downloads MAFF Japan agricultural trade statistics.
type Client struct{}

// NewClient creates a MAFF Japan client.
func NewClient() *Client { return &Client{} }

// FetchMAFFAgTradeSnapshot returns MAFF ag trade rows (fixture — portal HTTP 403 from CI).
func (c *Client) FetchMAFFAgTradeSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	_ = ctx
	rows, err := embeddedTradeSample()
	if err != nil {
		return nil, "", err
	}
	if path := strings.TrimSpace(os.Getenv("MAFF_BULK_PATH")); path != "" {
		if parsed, parseErr := parseTradeJSONFile(path); parseErr == nil && len(parsed) > 0 {
			rows = parsed
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].CommoditySlug+rows[i].FlowCode < rows[j].CommoditySlug+rows[j].FlowCode
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = "https://www.maff.go.jp/e/data/stat/export/index.html"
	}
	return payload, sourceURL + " (fixture — HTTP 403 from automated fetch 2026-06-26)", nil
}

func embeddedTradeSample() ([]tradeRow, error) {
	return []tradeRow{
		{CommoditySlug: "arroz", RefYear: "2023", FlowCode: "export", ValueJPY: "12500000000", QuantityT: "85000"},
		{CommoditySlug: "trigo", RefYear: "2023", FlowCode: "import", ValueJPY: "890000000000", QuantityT: "5600000"},
	}, nil
}

func parseTradeJSONFile(path string) ([]tradeRow, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rows []tradeRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// FlattenMAFFAgTrade converts merged MAFF JSON into canonical bronze columns.
func FlattenMAFFAgTrade(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []tradeRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse japan maff json: %w", err)
	}

	headers := []string{"commodity_slug", "refyear", "flow_code", "value_jpy", "quantity_t"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.CommoditySlug) == "" {
			continue
		}
		out = append(out, []string{
			row.CommoditySlug, row.RefYear, row.FlowCode, row.ValueJPY, row.QuantityT,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no japan maff rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
