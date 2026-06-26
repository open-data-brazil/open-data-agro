package usda

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const gatsAPIBaseURL = "https://apps.fas.usda.gov/OpenData/api/gats/trade"

type gatsRow struct {
	CommodityCode string `json:"commodity_code"`
	CommodityName string `json:"commodity_name"`
	PartnerCode   string `json:"partner_code"`
	PartnerName   string `json:"partner_name"`
	Flow          string `json:"flow"`
	Year          string `json:"year"`
	Value         string `json:"value"`
	Unit          string `json:"unit"`
}

type gatsAPIRecord struct {
	CommodityCode   int               `json:"commodityCode"`
	CommodityName   string            `json:"commodityName"`
	CountryCode     int               `json:"countryCode"`
	CountryName     string            `json:"countryName"`
	TradeFlow       string            `json:"tradeFlow"`
	UnitName        string            `json:"unitName"`
	YearlyValues    map[string]any    `json:"yearlyValues"`
	Values          map[string]any    `json:"values"`
	Year            int               `json:"year"`
	Value           float64           `json:"value"`
}

// FetchGATSSnapshot downloads U.S. agricultural trade rows via USDA FAS Open Data API.
func (c *Client) FetchGATSSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	apiKey := strings.TrimSpace(os.Getenv("USDA_FAS_API_KEY"))
	if apiKey == "" {
		return nil, "", fmt.Errorf("USDA_FAS_API_KEY is required for %s", entry.DatasetID)
	}

	commodities := entry.GATSCommodityCodes
	if len(commodities) == 0 {
		commodities = []string{"401", "801"}
	}
	partners := entry.GATSPartnerCodes
	if len(partners) == 0 {
		partners = []string{"5880", "1220"}
	}
	flows := entry.GATSFlows
	if len(flows) == 0 {
		flows = []string{"exports"}
	}
	years := entry.GATSYears
	if len(years) == 0 {
		years = defaultGATSYears(entry, fromDate)
	}
	dataValue := strings.TrimSpace(entry.GATSDataValue)
	if dataValue == "" {
		dataValue = "value"
	}

	var merged []gatsRow
	var firstURL string

	for _, commodity := range commodities {
		for _, partner := range partners {
			for _, flow := range flows {
				for _, year := range years {
					query := url.Values{
						"commodityCode": {strings.TrimSpace(commodity)},
						"partnerCode":   {strings.TrimSpace(partner)},
						"flow":          {strings.TrimSpace(flow)},
						"year":          {strconv.Itoa(year)},
						"dataValue":     {dataValue},
					}
					sourceURL := gatsAPIBaseURL + "?" + query.Encode()
					if firstURL == "" {
						firstURL = sourceURL
					}

					raw, err := c.downloadAPI(ctx, sourceURL, apiKey)
					if err != nil {
						return nil, "", fmt.Errorf("fetch gats %s/%s/%s/%d: %w", commodity, partner, flow, year, err)
					}

					rows, err := parseGATSResponse(raw, commodity, partner, flow, year, dataValue)
					if err != nil {
						return nil, "", err
					}
					merged = append(merged, rows...)
				}
			}
		}
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("usda gats returned no rows for %s", entry.DatasetID)
	}

	sort.Slice(merged, func(i, j int) bool {
		return gatsRowKey(merged[i]) < gatsRowKey(merged[j])
	})

	payload, err := json.Marshal(merged)
	if err != nil {
		return nil, "", err
	}
	return payload, firstURL, nil
}

func (c *Client) downloadAPI(ctx context.Context, sourceURL, apiKey string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Api-Key", apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d from %s", resp.StatusCode, sourceURL)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return nil, err
	}
	return body, nil
}

func parseGATSResponse(raw []byte, commodity, partner, flow string, year int, dataValue string) ([]gatsRow, error) {
	var records []gatsAPIRecord
	if err := json.Unmarshal(raw, &records); err == nil && len(records) > 0 {
		return flattenGATSRecords(records, commodity, partner, flow, year, dataValue), nil
	}

	var envelope struct {
		Data  []gatsAPIRecord `json:"data"`
		Rows  []gatsAPIRecord `json:"rows"`
		Value []gatsAPIRecord `json:"value"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("parse gats response: %w", err)
	}
	records = envelope.Data
	if len(records) == 0 {
		records = envelope.Rows
	}
	if len(records) == 0 {
		records = envelope.Value
	}
	if len(records) == 0 {
		return nil, nil
	}
	return flattenGATSRecords(records, commodity, partner, flow, year, dataValue), nil
}

func flattenGATSRecords(records []gatsAPIRecord, commodity, partner, flow string, year int, dataValue string) []gatsRow {
	yearStr := strconv.Itoa(year)
	var out []gatsRow
	for _, rec := range records {
		commodityCode := commodity
		if rec.CommodityCode > 0 {
			commodityCode = strconv.Itoa(rec.CommodityCode)
		}
		partnerCode := partner
		if rec.CountryCode > 0 {
			partnerCode = strconv.Itoa(rec.CountryCode)
		}
		flowName := flow
		if strings.TrimSpace(rec.TradeFlow) != "" {
			flowName = strings.TrimSpace(rec.TradeFlow)
		}
		unit := strings.TrimSpace(rec.UnitName)
		if unit == "" {
			unit = dataValue
		}

		if len(rec.YearlyValues) > 0 {
			for y, raw := range rec.YearlyValues {
				out = append(out, gatsRow{
					CommodityCode: commodityCode,
					CommodityName: strings.TrimSpace(rec.CommodityName),
					PartnerCode:   partnerCode,
					PartnerName:   strings.TrimSpace(rec.CountryName),
					Flow:          flowName,
					Year:          y,
					Value:         fmt.Sprint(raw),
					Unit:          unit,
				})
			}
			continue
		}
		if len(rec.Values) > 0 {
			for y, raw := range rec.Values {
				out = append(out, gatsRow{
					CommodityCode: commodityCode,
					CommodityName: strings.TrimSpace(rec.CommodityName),
					PartnerCode:   partnerCode,
					PartnerName:   strings.TrimSpace(rec.CountryName),
					Flow:          flowName,
					Year:          y,
					Value:         fmt.Sprint(raw),
					Unit:          unit,
				})
			}
			continue
		}

		rowYear := yearStr
		if rec.Year > 0 {
			rowYear = strconv.Itoa(rec.Year)
		}
		value := ""
		if rec.Value != 0 {
			value = strconv.FormatFloat(rec.Value, 'f', -1, 64)
		}
		out = append(out, gatsRow{
			CommodityCode: commodityCode,
			CommodityName: strings.TrimSpace(rec.CommodityName),
			PartnerCode:   partnerCode,
			PartnerName:   strings.TrimSpace(rec.CountryName),
			Flow:          flowName,
			Year:          rowYear,
			Value:         value,
			Unit:          unit,
		})
	}
	return out
}

func defaultGATSYears(entry catalog.RegistryEntry, fromDate string) []int {
	start := entry.PeriodStart
	if start == 0 {
		start = 2020
	}
	end := entry.PeriodEnd
	if end == 0 {
		end = 2024
	}
	if strings.TrimSpace(fromDate) != "" && len(fromDate) >= 4 {
		if y, err := strconv.Atoi(fromDate[:4]); err == nil && y > start {
			start = y
		}
	}
	var years []int
	for y := start; y <= end; y++ {
		years = append(years, y)
	}
	return years
}

func gatsRowKey(row gatsRow) string {
	return strings.Join([]string{row.CommodityCode, row.PartnerCode, row.Flow, row.Year}, "|")
}

// FlattenGATS converts merged GATS JSON into canonical bronze columns.
func FlattenGATS(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []gatsRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse gats json: %w", err)
	}

	headers := []string{
		"commodity_code",
		"commodity_name",
		"partner_code",
		"partner_name",
		"flow",
		"year",
		"value",
		"unit",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.Year) == "" || strings.TrimSpace(row.Value) == "" {
			continue
		}
		out = append(out, []string{
			row.CommodityCode,
			row.CommodityName,
			row.PartnerCode,
			row.PartnerName,
			row.Flow,
			row.Year,
			row.Value,
			row.Unit,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no gats rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
