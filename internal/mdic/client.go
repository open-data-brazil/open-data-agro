package mdic

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	comexAPIBase = "https://api-comexstat.mdic.gov.br"
	userAgent    = "open-data-agro/0.1.0 (+https://github.com/open-data-brazil/open-data-agro)"
)

// Client calls the MDIC Comex Stat API.
type Client struct {
	httpClient *http.Client
	maxRetries int
}

// NewClient creates an MDIC Comex Stat HTTP client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 180 * time.Second},
		maxRetries: 3,
	}
}

// PostJSON sends a POST request with a JSON body and returns the response body.
func (c *Client) PostJSON(ctx context.Context, path string, payload []byte) ([]byte, error) {
	url := comexAPIBase + path
	var lastErr error

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}

		body, err := c.postOnce(ctx, url, payload)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !isRetryable(err) {
			break
		}
	}

	return nil, fmt.Errorf("comex post failed after %d attempts: %w", c.maxRetries, lastErr)
}

func (c *Client) postOnce(ctx context.Context, url string, payload []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 128<<20))
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
