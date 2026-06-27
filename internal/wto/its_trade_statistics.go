package wto

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const (
	userAgent        = "open-data-agro/0.1.0 (+https://github.com/open-data-brazil/open-data-agro)"
	defaultWTOAPIURL = "https://api.wto.org/timeseries/v1/data"
)

type tradeRow struct {
	ReporterCode   string `json:"reporter_code"`
	ReporterName   string `json:"reporter_name"`
	PartnerCode    string `json:"partner_code"`
	PartnerName    string `json:"partner_name"`
	IndicatorCode  string `json:"indicator_code"`
	Period         string `json:"period"`
	ValueUSD       string `json:"value_usd"`
	FlowCode       string `json:"flow_code"`
}

// Client calls the WTO Stats API (subscription key required for live fetch).
type Client struct {
	httpClient *http.Client
	maxRetries int
}

// NewClient creates a WTO HTTP client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 120 * time.Second},
		maxRetries: 3,
	}
}

// FetchITSTradeSnapshot downloads WTO ITS trade statistics (API key or fixture fallback).
func (c *Client) FetchITSTradeSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	apiKey := strings.TrimSpace(os.Getenv("WTO_API_KEY"))
	if apiKey == "" {
		rows, err := embeddedITSSample()
		if err != nil {
			return nil, "", err
		}
		payload, err := json.Marshal(rows)
		if err != nil {
			return nil, "", err
		}
		return payload, "fixture:wto.its-trade-statistics (set WTO_API_KEY for live API)", nil
	}

	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = fmt.Sprintf("%s?i=HS_P_0070&r=840&p=000&ps=2023", defaultWTOAPIURL)
	}

	raw, err := c.download(ctx, sourceURL, apiKey)
	if err != nil {
		return nil, "", err
	}

	rows, err := parseWTOJSON(raw)
	if err != nil {
		return nil, "", err
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("wto its returned no rows for %s", entry.DatasetID)
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].ReporterCode+rows[i].Period < rows[j].ReporterCode+rows[j].Period
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	return payload, sourceURL, nil
}

func (c *Client) download(ctx context.Context, sourceURL, apiKey string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d from %s", resp.StatusCode, sourceURL)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return nil, err
	}
	return body, nil
}

func parseWTOJSON(raw []byte) ([]tradeRow, error) {
	var envelope struct {
		Data []tradeRow `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err == nil && len(envelope.Data) > 0 {
		return envelope.Data, nil
	}
	var rows []tradeRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, fmt.Errorf("parse wto json: %w", err)
	}
	return rows, nil
}

func embeddedITSSample() ([]tradeRow, error) {
	return []tradeRow{
		{ReporterCode: "840", ReporterName: "United States", PartnerCode: "000", PartnerName: "World", IndicatorCode: "HS_P_0070", Period: "2023", ValueUSD: "125000000", FlowCode: "X"},
		{ReporterCode: "76", ReporterName: "Brazil", PartnerCode: "000", PartnerName: "World", IndicatorCode: "HS_P_0070", Period: "2023", ValueUSD: "98000000", FlowCode: "X"},
	}, nil
}

// FlattenITSTrade converts merged WTO ITS JSON into canonical bronze columns.
func FlattenITSTrade(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []tradeRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse wto json: %w", err)
	}

	headers := []string{
		"reporter_code", "reporter_name", "partner_code", "partner_name",
		"indicator_code", "period", "value_usd", "flow_code",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.Period) == "" {
			continue
		}
		out = append(out, []string{
			row.ReporterCode, row.ReporterName, row.PartnerCode, row.PartnerName,
			row.IndicatorCode, row.Period, row.ValueUSD, row.FlowCode,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no wto rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
