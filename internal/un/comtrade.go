package un

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const (
	comtradeAPIBase    = "https://comtradeapi.un.org/data/v1/get/C/A/HS"
	comtradePreviewAPI = "https://comtradeapi.un.org/public/v1/preview/C/A/HS"
	comtradePortalURL  = "https://comtradeplus.un.org/"
)

// CommoditySlug maps HS chapter codes to canonical crop slugs.
var CommoditySlug = map[string]string{
	"1201": "soja",
	"1005": "milho",
	"1001": "trigo",
	"1006": "arroz",
	"1207": "oleo_vegetal",
}

type comtradeRow struct {
	ReporterCode   string `json:"reporter_code"`
	ReporterDesc   string `json:"reporter_desc"`
	PartnerCode    string `json:"partner_code"`
	PartnerDesc    string `json:"partner_desc"`
	FlowCode       string `json:"flow_code"`
	FlowDesc       string `json:"flow_desc"`
	Period         string `json:"period"`
	HSCode         string `json:"hs_code"`
	CommoditySlug  string `json:"commodity_slug"`
	TradeValueUSD  string `json:"trade_value_usd"`
	NetWeightKg    string `json:"netweight_kg"`
	Qty            string `json:"qty"`
	QtyUnitAbbr    string `json:"qty_unit_abbr"`
}

type comtradeAPIResponse struct {
	Count int             `json:"count"`
	Data  []comtradeAPIData `json:"data"`
}

type comtradeAPIData struct {
	ReporterCode  int     `json:"reporterCode"`
	ReporterDesc  string  `json:"reporterDesc"`
	PartnerCode   int     `json:"partnerCode"`
	PartnerDesc   string  `json:"partnerDesc"`
	FlowCode      string  `json:"flowCode"`
	FlowDesc      string  `json:"flowDesc"`
	Period        string  `json:"period"`
	CmdCode       string  `json:"cmdCode"`
	PrimaryValue  float64 `json:"primaryValue"`
	NetWgt        float64 `json:"netWgt"`
	Qty           float64 `json:"qty"`
	QtyUnitAbbr   string  `json:"qtyUnitAbbr"`
}

// FetchComtradeSnapshot downloads bilateral ag trade rows via UN Comtrade API.
func (c *Client) FetchComtradeSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	reporter := strings.TrimSpace(entry.ComtradeReporterCode)
	if reporter == "" {
		reporter = "76"
	}

	cmdCodes := entry.ComtradeCmdCodes
	if len(cmdCodes) == 0 {
		cmdCodes = []string{"1201", "1005", "1001", "1006"}
	}

	flowCodes := entry.ComtradeFlowCodes
	if len(flowCodes) == 0 {
		flowCodes = []string{"X", "M"}
	}

	startYear, endYear := resolveComtradeYearRange(entry, fromDate)
	apiKey := strings.TrimSpace(os.Getenv("COMTRADE_SUBSCRIPTION_KEY"))

	merged := make(map[string]comtradeRow)
	requests := 0

	for year := startYear; year <= endYear; year++ {
		period := strconv.Itoa(year)
		for _, flow := range flowCodes {
			requestURL, err := buildComtradeURL(reporter, period, flow, cmdCodes, apiKey)
			if err != nil {
				return nil, "", err
			}

			raw, err := c.Download(ctx, requestURL)
			if err != nil {
				continue
			}
			requests++

			rows, err := parseComtradeResponse(raw)
			if err != nil {
				continue
			}
			for _, row := range rows {
				key := strings.Join([]string{
					row.ReporterCode, row.PartnerCode, row.FlowCode, row.Period, row.HSCode,
				}, "|")
				merged[key] = row
			}
		}
	}

	if len(merged) == 0 {
		if apiKey == "" {
			return nil, "", fmt.Errorf("COMTRADE_SUBSCRIPTION_KEY is required for %s (free tier: 500 calls/day, 100k records/call)", entry.DatasetID)
		}
		return nil, "", fmt.Errorf("un comtrade returned no rows for %s", entry.DatasetID)
	}

	keys := make([]string, 0, len(merged))
	for key := range merged {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	out := make([]comtradeRow, 0, len(keys))
	for _, key := range keys {
		out = append(out, merged[key])
	}

	payload, err := json.Marshal(out)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s?reporter=%s&years=%d-%d&requests=%d", comtradePortalURL, reporter, startYear, endYear, requests)
	return payload, sourceURL, nil
}

func buildComtradeURL(reporter, period, flow string, cmdCodes []string, apiKey string) (string, error) {
	base := comtradeAPIBase
	if apiKey == "" {
		base = comtradePreviewAPI
	}

	values := url.Values{}
	values.Set("reporterCode", reporter)
	values.Set("partnerCode", "0")
	values.Set("period", period)
	values.Set("flowCode", flow)
	values.Set("cmdCode", strings.Join(cmdCodes, ","))
	values.Set("maxRecords", "100000")
	if apiKey != "" {
		values.Set("subscription-key", apiKey)
	}

	parsed, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func parseComtradeResponse(raw []byte) ([]comtradeRow, error) {
	var payload comtradeAPIResponse
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}

	rows := make([]comtradeRow, 0, len(payload.Data))
	for _, item := range payload.Data {
		if strings.TrimSpace(item.Period) == "" || strings.TrimSpace(item.CmdCode) == "" {
			continue
		}
		rows = append(rows, comtradeRow{
			ReporterCode:  strconv.Itoa(item.ReporterCode),
			ReporterDesc:  strings.TrimSpace(item.ReporterDesc),
			PartnerCode:   strconv.Itoa(item.PartnerCode),
			PartnerDesc:   strings.TrimSpace(item.PartnerDesc),
			FlowCode:      strings.TrimSpace(item.FlowCode),
			FlowDesc:      strings.TrimSpace(item.FlowDesc),
			Period:        strings.TrimSpace(item.Period),
			HSCode:        strings.TrimSpace(item.CmdCode),
			CommoditySlug: commoditySlug(item.CmdCode),
			TradeValueUSD: strconv.FormatFloat(item.PrimaryValue, 'f', -1, 64),
			NetWeightKg:   strconv.FormatFloat(item.NetWgt, 'f', -1, 64),
			Qty:           strconv.FormatFloat(item.Qty, 'f', -1, 64),
			QtyUnitAbbr:   strings.TrimSpace(item.QtyUnitAbbr),
		})
	}
	return rows, nil
}

func commoditySlug(hsCode string) string {
	code := strings.TrimSpace(hsCode)
	if slug, ok := CommoditySlug[code]; ok {
		return slug
	}
	if len(code) >= 4 {
		if slug, ok := CommoditySlug[code[:4]]; ok {
			return slug
		}
	}
	return strings.ToLower(code)
}

func resolveComtradeYearRange(entry catalog.RegistryEntry, fromDate string) (int, int) {
	startYear := entry.PeriodStart
	if startYear == 0 {
		startYear = 2010
	}

	fromDate = strings.TrimSpace(fromDate)
	if fromDate != "" {
		if len(fromDate) >= 4 {
			if year, err := strconv.Atoi(fromDate[:4]); err == nil && year > startYear {
				startYear = year
			}
		}
	}

	endYear := time.Now().UTC().Year()
	if entry.PeriodEnd > 0 && entry.PeriodEnd < endYear {
		endYear = entry.PeriodEnd
	}
	return startYear, endYear
}

// FlattenComtrade converts merged Comtrade JSON into canonical bronze columns.
func FlattenComtrade(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []comtradeRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse un comtrade json: %w", err)
	}

	headers := []string{
		"reporter_code",
		"reporter_desc",
		"partner_code",
		"partner_desc",
		"flow_code",
		"flow_desc",
		"period",
		"hs_code",
		"commodity_slug",
		"trade_value_usd",
		"netweight_kg",
		"qty",
		"qty_unit_abbr",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.Period) == "" || strings.TrimSpace(row.HSCode) == "" {
			continue
		}
		out = append(out, []string{
			row.ReporterCode,
			row.ReporterDesc,
			row.PartnerCode,
			row.PartnerDesc,
			row.FlowCode,
			row.FlowDesc,
			row.Period,
			row.HSCode,
			row.CommoditySlug,
			row.TradeValueUSD,
			row.NetWeightKg,
			row.Qty,
			row.QtyUnitAbbr,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no un comtrade rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

// ResolveURL validates the catalog base URL for a UN Comtrade dataset.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return comtradePortalURL, nil
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid source_url for %s: %w", entry.DatasetID, err)
	}
	if parsed.Scheme != "https" {
		return "", fmt.Errorf("source_url for %s must use https", entry.DatasetID)
	}
	host := strings.ToLower(parsed.Host)
	if !strings.Contains(host, "un.org") {
		return "", fmt.Errorf("source_url for %s must be on un.org", entry.DatasetID)
	}
	return parsed.String(), nil
}
