package argentina

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const userAgent = "open-data-agro/0.1.0 (+https://github.com/open-data-brazil/open-data-agro)"

const defaultCambioBaseURL = "https://api.bcra.gob.ar/estadisticascambiarias/v1.0/Cotizaciones"

type cambioRow struct {
	CurrencyCode string `json:"currency_code"`
	CurrencyName string `json:"currency_name"`
	RefDate      string `json:"refdate"`
	ExchangeRate string `json:"exchange_rate"`
	RateType     string `json:"rate_type"`
}

type cambioAPIResponse struct {
	Status  int `json:"status"`
	Results []struct {
		Fecha   string `json:"fecha"`
		Detalle []struct {
			CodigoMoneda    string  `json:"codigoMoneda"`
			Descripcion     string  `json:"descripcion"`
			TipoPase        float64 `json:"tipoPase"`
			TipoCotizacion  float64 `json:"tipoCotizacion"`
		} `json:"detalle"`
	} `json:"results"`
}

// Client downloads official statistics from BCRA.
type Client struct {
	httpClient *http.Client
	maxRetries int
}

// NewClient creates a BCRA HTTP client with retry and timeout defaults.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 120 * time.Second},
		maxRetries: 3,
	}
}

// FetchCambioSnapshot downloads BCRA official exchange-rate series.
func (c *Client) FetchCambioSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	currency := strings.TrimSpace(entry.ArgentinaCurrencyCode)
	if currency == "" {
		currency = "USD"
	}

	fechaDesde := strings.TrimSpace(entry.StartDate)
	if fechaDesde == "" {
		fechaDesde = "2024-01-01"
	}
	if strings.TrimSpace(fromDate) != "" {
		fechaDesde = strings.TrimSpace(fromDate)
	}
	fechaHasta := time.Now().Format("2006-01-02")

	query := url.Values{
		"fechaDesde": {fechaDesde},
		"fechaHasta": {fechaHasta},
	}
	sourceURL := fmt.Sprintf("%s/%s?%s", defaultCambioBaseURL, url.PathEscape(currency), query.Encode())

	raw, err := c.download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	rows, err := parseCambioJSON(raw, currency)
	if err != nil {
		return nil, "", err
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("bcra cambio returned no rows for %s", entry.DatasetID)
	}

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	return payload, sourceURL, nil
}

func parseCambioJSON(raw []byte, currency string) ([]cambioRow, error) {
	var payload cambioAPIResponse
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("parse bcra cambio json: %w", err)
	}

	var rows []cambioRow
	for _, day := range payload.Results {
		for _, detail := range day.Detalle {
			if currency != "" && !strings.EqualFold(strings.TrimSpace(detail.CodigoMoneda), currency) {
				continue
			}
			rate := detail.TipoCotizacion
			rateType := "tipo_cotizacion"
			if rate == 0 && detail.TipoPase != 0 {
				rate = detail.TipoPase
				rateType = "tipo_pase"
			}
			rows = append(rows, cambioRow{
				CurrencyCode: strings.TrimSpace(detail.CodigoMoneda),
				CurrencyName: strings.TrimSpace(detail.Descripcion),
				RefDate:      strings.TrimSpace(day.Fecha),
				ExchangeRate: fmt.Sprintf("%.8f", rate),
				RateType:     rateType,
			})
		}
	}
	return rows, nil
}

func (c *Client) download(ctx context.Context, sourceURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
		body, err := c.downloadOnce(ctx, sourceURL)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !isRetryable(err) {
			break
		}
	}
	return nil, fmt.Errorf("bcra download failed after %d attempts: %w", c.maxRetries, lastErr)
}

func (c *Client) downloadOnce(ctx context.Context, sourceURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

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

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "status 5")
}

// FlattenCambio converts merged BCRA JSON into canonical bronze columns.
func FlattenCambio(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []cambioRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse bcra cambio json: %w", err)
	}

	headers := []string{
		"currency_code",
		"currency_name",
		"refdate",
		"exchange_rate",
		"rate_type",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RefDate) == "" || strings.TrimSpace(row.ExchangeRate) == "" {
			continue
		}
		out = append(out, []string{
			row.CurrencyCode,
			row.CurrencyName,
			row.RefDate,
			row.ExchangeRate,
			row.RateType,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no bcra cambio rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

// ResolveURL returns the BCRA exchange-rate API URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	currency := strings.TrimSpace(entry.ArgentinaCurrencyCode)
	if currency == "" {
		currency = "USD"
	}
	fechaDesde := strings.TrimSpace(entry.StartDate)
	if fechaDesde == "" {
		fechaDesde = "2024-01-01"
	}
	query := url.Values{
		"fechaDesde": {fechaDesde},
		"fechaHasta": {time.Now().Format("2006-01-02")},
	}
	return fmt.Sprintf("%s/%s?%s", defaultCambioBaseURL, url.PathEscape(currency), query.Encode()), nil
}
