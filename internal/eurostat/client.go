package eurostat

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const userAgent = "open-data-agro/0.1.0 (+https://github.com/open-data-brazil/open-data-agro)"

const defaultAPIBaseURL = "https://ec.europa.eu/eurostat/api/dissemination/statistics/1.0/data"

// Client downloads datasets from the Eurostat JSON API.
type Client struct {
	httpClient *http.Client
	maxRetries int
}

// NewClient creates a Eurostat HTTP client with retry and timeout defaults.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 180 * time.Second},
		maxRetries: 3,
	}
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
	return nil, fmt.Errorf("eurostat download failed after %d attempts: %w", c.maxRetries, lastErr)
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

	body, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
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

func buildDatasetURL(datasetCode, geo string, products []string, sinceYear int) string {
	base := defaultAPIBaseURL + "/" + url.PathEscape(strings.TrimSpace(datasetCode))
	values := url.Values{}
	values.Set("format", "JSON")
	values.Set("lang", "en")
	if geo != "" {
		values.Set("geo", geo)
	}
	if sinceYear > 0 {
		values.Set("sinceTimePeriod", strconvItoa(sinceYear))
	}
	for _, product := range products {
		if strings.TrimSpace(product) != "" {
			values.Add("product", strings.TrimSpace(product))
		}
	}
	return base + "?" + values.Encode()
}

func strconvItoa(v int) string {
	return fmt.Sprintf("%d", v)
}
