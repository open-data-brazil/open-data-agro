package cftc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const (
	defaultCOTSocrataURL = "https://publicreporting.cftc.gov/resource/72hh-3qpy.json"
	defaultCOTLimit      = 50000
)

var commoditySlug = map[string]string{
	"CORN":     "milho",
	"SOYBEANS": "soja",
	"WHEAT":    "trigo",
	"RICE":     "arroz",
	"COTTON":   "algodao",
	"SUGAR":    "acucar",
	"COFFEE":   "cafe",
	"COCOA":    "cacau",
	"OATS":     "aveia",
	"SOYBEAN OIL": "oleo_soja",
}

type cotRow struct {
	ReportDate           string `json:"report_date"`
	CommodityName        string `json:"commodity_name"`
	CommoditySlug        string `json:"commodity_slug"`
	MarketName           string `json:"market_name"`
	OpenInterestAll      string `json:"open_interest_all"`
	MMoneyLong           string `json:"m_money_positions_long_all"`
	MMoneyShort          string `json:"m_money_positions_short_all"`
	ProdMercLong         string `json:"prod_merc_positions_long"`
	ProdMercShort        string `json:"prod_merc_positions_short"`
	CommodityGroup       string `json:"commodity_group_name"`
	FutOnlyOrCombined    string `json:"futonly_or_combined"`
}

type cotAPIRow struct {
	ReportDateAsYYYYMMD  string `json:"report_date_as_yyyy_mm_dd"`
	CommodityName        string `json:"commodity_name"`
	MarketAndExchange    string `json:"market_and_exchange_names"`
	OpenInterestAll      string `json:"open_interest_all"`
	MMoneyLong           string `json:"m_money_positions_long_all"`
	MMoneyShort          string `json:"m_money_positions_short_all"`
	ProdMercLong         string `json:"prod_merc_positions_long"`
	ProdMercShort        string `json:"prod_merc_positions_short"`
	CommodityGroup       string `json:"commodity_group_name"`
	FutOnlyOrCombined    string `json:"futonly_or_combined"`
}

// FetchCOTAgriculturalSnapshot downloads disaggregated COT futures-only agricultural positions.
func (c *Client) FetchCOTAgriculturalSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	sourceURL, err := buildCOTURL(entry)
	if err != nil {
		return nil, "", err
	}

	raw, err := c.download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	rows, err := parseCOTJSON(raw)
	if err != nil {
		return nil, "", err
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("cftc cot returned no agricultural rows for %s", entry.DatasetID)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].CommoditySlug != rows[j].CommoditySlug {
			return rows[i].CommoditySlug < rows[j].CommoditySlug
		}
		return rows[i].ReportDate < rows[j].ReportDate
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	return payload, sourceURL, nil
}

func buildCOTURL(entry catalog.RegistryEntry) (string, error) {
	base := strings.TrimSpace(entry.SourceURL)
	if base == "" {
		base = defaultCOTSocrataURL
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("parse cftc source url: %w", err)
	}
	q := parsed.Query()
	if q.Get("$where") == "" {
		q.Set("$where", "commodity_group_name='AGRICULTURE' AND futonly_or_combined='FutOnly'")
	}
	if q.Get("$limit") == "" {
		q.Set("$limit", fmt.Sprintf("%d", defaultCOTLimit))
	}
	if q.Get("$order") == "" {
		q.Set("$order", "report_date_as_yyyy_mm_dd DESC")
	}
	parsed.RawQuery = q.Encode()
	return parsed.String(), nil
}

func parseCOTJSON(raw []byte) ([]cotRow, error) {
	var apiRows []cotAPIRow
	if err := json.Unmarshal(raw, &apiRows); err != nil {
		return nil, fmt.Errorf("parse cftc cot json: %w", err)
	}

	out := make([]cotRow, 0, len(apiRows))
	for _, row := range apiRows {
		if !strings.EqualFold(strings.TrimSpace(row.CommodityGroup), "AGRICULTURE") {
			continue
		}
		reportDate := normalizeCOTDate(row.ReportDateAsYYYYMMD)
		if reportDate == "" {
			continue
		}
		name := strings.TrimSpace(row.CommodityName)
		slug := slugForCommodity(name)
		out = append(out, cotRow{
			ReportDate:        reportDate,
			CommodityName:     name,
			CommoditySlug:     slug,
			MarketName:        strings.TrimSpace(row.MarketAndExchange),
			OpenInterestAll:   strings.TrimSpace(row.OpenInterestAll),
			MMoneyLong:        strings.TrimSpace(row.MMoneyLong),
			MMoneyShort:       strings.TrimSpace(row.MMoneyShort),
			ProdMercLong:      strings.TrimSpace(row.ProdMercLong),
			ProdMercShort:     strings.TrimSpace(row.ProdMercShort),
			CommodityGroup:    strings.TrimSpace(row.CommodityGroup),
			FutOnlyOrCombined: strings.TrimSpace(row.FutOnlyOrCombined),
		})
	}
	return out, nil
}

func normalizeCOTDate(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.Contains(raw, "T") {
		raw = strings.Split(raw, "T")[0]
	}
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return t.Format("2006-01-02")
	}
	return raw
}

func slugForCommodity(name string) string {
	upper := strings.ToUpper(strings.TrimSpace(name))
	if slug, ok := commoditySlug[upper]; ok {
		return slug
	}
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "_"))
}

// FlattenCOTAgricultural converts merged CFTC COT JSON into canonical bronze columns.
func FlattenCOTAgricultural(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []cotRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse cftc cot json: %w", err)
	}

	headers := []string{
		"report_date", "commodity_name", "commodity_slug", "market_name",
		"open_interest_all", "m_money_long", "m_money_short",
		"prod_merc_long", "prod_merc_short", "commodity_group", "futonly_or_combined",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.ReportDate) == "" {
			continue
		}
		out = append(out, []string{
			row.ReportDate,
			row.CommodityName,
			row.CommoditySlug,
			row.MarketName,
			row.OpenInterestAll,
			row.MMoneyLong,
			row.MMoneyShort,
			row.ProdMercLong,
			row.ProdMercShort,
			row.CommodityGroup,
			row.FutOnlyOrCombined,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no cftc cot rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
