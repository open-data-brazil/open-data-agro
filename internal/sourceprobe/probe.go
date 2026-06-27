package sourceprobe

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	UserAgent       = "open-data-agro/0.1.0 (+https://github.com/open-data-brazil/open-data-agro)"
	DefaultSampleSz = 4096
	DefaultTimeout  = 60 * time.Second
	DefaultRetries  = 3
)

// Client probes official source URLs with GET sample downloads (same transport as ingest).
type Client struct {
	httpClient *http.Client
	maxRetries int
	sampleSize int
}

// NewClient creates a probe client with ingest-compatible defaults.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: DefaultTimeout},
		maxRetries: DefaultRetries,
		sampleSize: DefaultSampleSz,
	}
}

// SampleResult holds bytes and HTTP metadata from a probe GET.
type SampleResult struct {
	HTTPStatus   int
	SampleBytes  int
	SampleSHA256 string
	LastModified string
	ContentType  string
}

// ProbeURL performs a GET with a limited body read (sample only).
func (c *Client) ProbeURL(ctx context.Context, url string) (*SampleResult, error) {
	return c.ProbeURLWithHeaders(ctx, url, map[string]string{"Accept": "*/*"})
}

// ProbeURLWithHeaders performs a GET sample probe with extra request headers (ingest-compatible).
func (c *Client) ProbeURLWithHeaders(ctx context.Context, url string, headers map[string]string) (*SampleResult, error) {
	var lastErr error

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}

		result, err := c.probeOnceWithHeaders(ctx, url, headers)
		if err == nil {
			return result, nil
		}
		lastErr = err
		if !isRetryable(err) {
			break
		}
	}

	return nil, fmt.Errorf("probe failed after %d attempts: %w", c.maxRetries, lastErr)
}

func (c *Client) probeOnceWithHeaders(ctx context.Context, url string, headers map[string]string) (*SampleResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "*/*")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(c.sampleSize)))
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response body for %s", url)
	}

	sum := sha256.Sum256(body)
	return &SampleResult{
		HTTPStatus:   resp.StatusCode,
		SampleBytes:  len(body),
		SampleSHA256: hex.EncodeToString(sum[:]),
		LastModified: strings.TrimSpace(resp.Header.Get("Last-Modified")),
		ContentType:  strings.TrimSpace(resp.Header.Get("Content-Type")),
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
