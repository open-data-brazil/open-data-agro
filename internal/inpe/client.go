package inpe

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

const (
	userAgent        = "open-data-agro/0.1.0 (+https://github.com/open-data-brazil/open-data-agro)"
	defaultWFSBase   = "https://terrabrasilis.dpi.inpe.br/geoserver/deter-amz/wfs"
	defaultTypeName  = "deter-amz:deter_amz"
	defaultMaxAlerts = 5000
)

// DownloadResult holds a fetched INPE response and HTTP metadata.
type DownloadResult struct {
	Body          []byte
	ContentType   string
	LastModified  string
	ContentLength int64
	SourceURL     string
}

// Client downloads DETER alerts from TerraBrasilis WFS.
type Client struct {
	httpClient *http.Client
	maxRetries int
}

// NewClient creates an INPE HTTP client with retry and timeout defaults.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 600 * time.Second},
		maxRetries: 3,
	}
}

// FetchDETERSnapshot downloads recent DETER alerts as GeoJSON features.
func (c *Client) FetchDETERSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	requestURL, err := buildWFSURL(entry)
	if err != nil {
		return nil, "", err
	}

	result, err := c.Download(ctx, requestURL)
	if err != nil {
		return nil, "", fmt.Errorf("inpe deter fetch: %w", err)
	}

	rows, err := parseFeatureCollection(result.Body)
	if err != nil {
		return nil, "", err
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("inpe deter returned no alert rows for %s", entry.DatasetID)
	}

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}

	return payload, requestURL, nil
}

// Download fetches bytes from a validated INPE portal URL.
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
	req.Header.Set("Accept", "application/json")

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

type geoFeatureCollection struct {
	Features []geoFeature `json:"features"`
}

type geoFeature struct {
	Properties map[string]any `json:"properties"`
}

func parseFeatureCollection(raw []byte) ([]map[string]any, error) {
	var payload geoFeatureCollection
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("parse deter geojson: %w", err)
	}
	rows := make([]map[string]any, 0, len(payload.Features))
	for _, feature := range payload.Features {
		if len(feature.Properties) == 0 {
			continue
		}
		rows = append(rows, feature.Properties)
	}
	return rows, nil
}

func parseDETERRows(raw []byte) ([]map[string]any, error) {
	var merged []map[string]any
	if err := json.Unmarshal(raw, &merged); err == nil && len(merged) > 0 {
		return merged, nil
	}

	rows, err := parseFeatureCollection(raw)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func buildWFSURL(entry catalog.RegistryEntry) (string, error) {
	base := strings.TrimSpace(entry.SourceURL)
	if base == "" {
		base = defaultWFSBase
	}

	count := defaultMaxAlerts
	if entry.PeriodStart > 0 && entry.PeriodStart < count {
		count = entry.PeriodStart
	}

	params := url.Values{}
	params.Set("service", "WFS")
	params.Set("version", "2.0.0")
	params.Set("request", "GetFeature")
	params.Set("typeName", defaultTypeName)
	params.Set("count", fmt.Sprintf("%d", count))
	params.Set("outputFormat", "application/json")
	params.Set("propertyName", "view_date,areauckm,uf,classname,municipality,publish_month")

	if strings.Contains(base, "?") {
		return base + "&" + params.Encode(), nil
	}
	return base + "?" + params.Encode(), nil
}

// FlattenDETER converts merged DETER alert JSON rows into bronze columns.
func FlattenDETER(raw []byte) ([]string, [][]string, error) {
	rows, err := parseDETERRows(raw)
	if err != nil {
		return nil, nil, err
	}

	headers := []string{
		"view_date",
		"class_name",
		"uf",
		"municipality",
		"area_uc_km",
		"publish_month",
	}

	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			stringField(row["view_date"]),
			stringField(row["classname"]),
			stringField(row["uf"]),
			stringField(row["municipality"]),
			formatNumber(row["areauckm"]),
			stringField(row["publish_month"]),
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no deter rows to flatten")
	}
	return headers, out, nil
}

func stringField(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case float64:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.4f", typed), "0"), ".")
	case nil:
		return ""
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}

func formatNumber(value any) string {
	switch typed := value.(type) {
	case float64:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", typed), "0"), ".")
	case int:
		return fmt.Sprintf("%d", typed)
	default:
		return stringField(value)
	}
}

// ResolveURL validates the catalog base URL for a DETER dataset.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return "", fmt.Errorf("dataset %s has no source_url", entry.DatasetID)
	}
	if !strings.Contains(strings.ToLower(raw), "inpe.br") {
		return "", fmt.Errorf("source_url for %s must be on inpe.br", entry.DatasetID)
	}
	return raw, nil
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
