package sourceprobe

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPRequest describes a probe HTTP call (GET or POST sample).
type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

// ProbeHTTP executes a sample probe with the given HTTP request.
func (c *Client) ProbeHTTP(ctx context.Context, req HTTPRequest) (*SampleResult, error) {
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	if method == "" {
		method = http.MethodGet
	}

	var lastErr error
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}

		result, err := c.probeHTTPOnce(ctx, method, req.URL, req.Headers, req.Body)
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

func (c *Client) probeHTTPOnce(ctx context.Context, method, url string, headers map[string]string, body string) (*SampleResult, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("User-Agent", UserAgent)
	for key, value := range headers {
		httpReq.Header.Set(key, value)
	}
	if httpReq.Header.Get("Accept") == "" {
		httpReq.Header.Set("Accept", "*/*")
	}
	if body != "" && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	sample, err := io.ReadAll(io.LimitReader(resp.Body, int64(c.sampleSize)))
	if err != nil {
		return nil, err
	}
	if len(sample) == 0 {
		return nil, fmt.Errorf("empty response body for %s", url)
	}

	hash := sha256.Sum256(sample)
	return &SampleResult{
		HTTPStatus:   resp.StatusCode,
		SampleBytes:  len(sample),
		SampleSHA256: hex.EncodeToString(hash[:]),
		LastModified: strings.TrimSpace(resp.Header.Get("Last-Modified")),
		ContentType:  strings.TrimSpace(resp.Header.Get("Content-Type")),
	}, nil
}
