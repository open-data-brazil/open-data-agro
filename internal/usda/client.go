package usda

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	userAgent      = "open-data-agro/0.1.0 (+https://github.com/open-data-brazil/open-data-agro)"
	psdSOAPEndpoint = "https://apps.fas.usda.gov/PSDExternalAPIService/svcPSD_AMIS.asmx"
)

// Client calls the USDA FAS PSD AMIS SOAP web service (no API key required).
type Client struct {
	httpClient *http.Client
	maxRetries int
	minGap     time.Duration
}

// NewClient creates a USDA PSD HTTP client with retry and timeout defaults.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 180 * time.Second},
		maxRetries: 3,
		minGap:     400 * time.Millisecond,
	}
}

// PostSOAP sends a SOAP request and returns the response body.
func (c *Client) PostSOAP(ctx context.Context, soapAction, envelope string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}

		body, err := c.postOnce(ctx, soapAction, envelope)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !isRetryable(err) {
			break
		}
	}

	return nil, fmt.Errorf("usda soap failed after %d attempts: %w", c.maxRetries, lastErr)
}

func (c *Client) postOnce(ctx context.Context, soapAction, envelope string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, psdSOAPEndpoint, bytes.NewReader([]byte(envelope)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", soapAction)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d from usda psd soap", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 256<<20))
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
	case <-time.After(c.minGap):
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
