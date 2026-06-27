package abiove

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const userAgent = "open-data-agro/0.1.0 (+https://github.com/open-data-brazil/open-data-agro)"

// DownloadResult holds a fetched Abiove statistics file.
type DownloadResult struct {
	Body          []byte
	ContentType   string
	LastModified  string
	ContentLength int64
	SourceURL     string
}

// Client downloads Excel reports from abiove.org.br.
type Client struct {
	httpClient *http.Client
	maxRetries int
}

// NewClient creates an Abiove HTTP client with retry defaults.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 120 * time.Second},
		maxRetries: 3,
	}
}

// Download fetches bytes from a validated Abiove URL.
func (c *Client) Download(ctx context.Context, sourceURL string) (*DownloadResult, error) {
	var lastErr error

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}

		result, err := c.downloadOnce(ctx, sourceURL)
		if err == nil {
			return result, nil
		}
		lastErr = err
		if !isRetryable(err) {
			break
		}
	}

	return nil, fmt.Errorf("download failed after %d attempts: %w", c.maxRetries, lastErr)
}

func (c *Client) downloadOnce(ctx context.Context, sourceURL string) (*DownloadResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, sourceURL)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
	if err != nil {
		return nil, err
	}

	return &DownloadResult{
		Body:          body,
		ContentType:   strings.TrimSpace(resp.Header.Get("Content-Type")),
		LastModified:  strings.TrimSpace(resp.Header.Get("Last-Modified")),
		ContentLength: int64(len(body)),
		SourceURL:     sourceURL,
	}, nil
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
