package ons

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const (
	userAgent              = "open-data-agro/0.1.0 (+https://github.com/open-data-brazil/open-data-agro)"
	defaultCKANPackageShow = "https://dados.ons.org.br/api/3/action/package_show"
)

var yearInNamePattern = regexp.MustCompile(`\b(20\d{2})\b`)

type ckanPackageResponse struct {
	Success bool `json:"success"`
	Result  struct {
		Resources []ckanResource `json:"resources"`
	} `json:"result"`
}

type ckanResource struct {
	Name         string `json:"name"`
	Format       string `json:"format"`
	URL          string `json:"url"`
	LastModified string `json:"last_modified"`
}

// DownloadResult holds a fetched ONS file and HTTP metadata.
type DownloadResult struct {
	Body          []byte
	ContentType   string
	LastModified  string
	ContentLength int64
	SourceURL     string
}

// Client downloads files from dados.ons.org.br.
type Client struct {
	httpClient *http.Client
	maxRetries int
}

// NewClient creates an ONS HTTP client with retry and timeout defaults.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 300 * time.Second},
		maxRetries: 3,
	}
}

// Download fetches bytes from a validated ONS portal URL.
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

	body, err := io.ReadAll(io.LimitReader(resp.Body, 128<<20))
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

// ResolveURL returns the latest annual CSV for an ONS CKAN package.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	packageID := strings.TrimSpace(entry.CKANPackageID)
	if packageID == "" {
		return resolveDirectURL(entry)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return resolveLatestCSV(ctx, packageID)
}

func resolveDirectURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return "", fmt.Errorf("dataset %s has no source_url or ckan_package_id", entry.DatasetID)
	}
	lower := strings.ToLower(raw)
	if !strings.Contains(lower, "ons.org.br") && !strings.Contains(lower, "amazonaws.com") {
		return "", fmt.Errorf("source_url for %s must be on ons.org.br or ONS S3", entry.DatasetID)
	}
	return raw, nil
}

func resolveLatestCSV(ctx context.Context, packageID string) (string, error) {
	url := fmt.Sprintf("%s?id=%s", defaultCKANPackageShow, packageID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return "", err
	}

	var payload ckanPackageResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("parse ons ckan response: %w", err)
	}
	if !payload.Success {
		return "", fmt.Errorf("ons ckan package_show failed for %s", packageID)
	}

	var matches []ckanResource
	for _, res := range payload.Result.Resources {
		if strings.ToUpper(strings.TrimSpace(res.Format)) != "CSV" {
			continue
		}
		if strings.TrimSpace(res.URL) == "" {
			continue
		}
		matches = append(matches, res)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no CSV resource in ons package %s", packageID)
	}

	sort.Slice(matches, func(i, j int) bool {
		yearI, yearJ := yearInNameSortKey(matches[i].Name), yearInNameSortKey(matches[j].Name)
		if yearI != yearJ {
			return yearI > yearJ
		}
		return matches[i].LastModified > matches[j].LastModified
	})

	return matches[0].URL, nil
}

func yearInNameSortKey(name string) int {
	match := yearInNamePattern.FindStringSubmatch(name)
	if len(match) < 2 {
		return -1
	}
	year, err := strconv.Atoi(match[1])
	if err != nil {
		return -1
	}
	return year
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
