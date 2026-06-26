package antt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const ckanAPIBase = "https://dados.antt.gov.br/api/3/action/package_show"

type ckanPackageResponse struct {
	Success bool `json:"success"`
	Result  struct {
		Resources []ckanResource `json:"resources"`
	} `json:"result"`
}

type ckanResource struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Format       string `json:"format"`
	URL          string `json:"url"`
	LastModified string `json:"last_modified"`
}

// ResolveURL returns the latest CKAN resource download URL for an ANTT catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	packageID := strings.TrimSpace(entry.CKANPackageID)
	if packageID == "" {
		return resolveDirectURL(entry)
	}

	format := strings.TrimSpace(entry.CKANResourceFormat)
	if format == "" {
		format = "CSV"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return ResolveLatestCKANResourceURL(ctx, packageID, format)
}

func resolveDirectURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return "", fmt.Errorf("dataset %s has no source_url or ckan_package_id", entry.DatasetID)
	}
	if !strings.Contains(strings.ToLower(raw), "gov.br") {
		return "", fmt.Errorf("source_url for %s must be on gov.br", entry.DatasetID)
	}
	return raw, nil
}

// ResolveLatestCKANResourceURL picks the newest resource matching format from a CKAN package.
func ResolveLatestCKANResourceURL(ctx context.Context, packageID, format string) (string, error) {
	url := fmt.Sprintf("%s?id=%s", ckanAPIBase, packageID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ckan package_show: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ckan status %d for package %s", resp.StatusCode, packageID)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return "", err
	}

	var payload ckanPackageResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("parse ckan response: %w", err)
	}
	if !payload.Success {
		return "", fmt.Errorf("ckan package_show failed for %s", packageID)
	}

	wantFormat := strings.ToUpper(strings.TrimSpace(format))
	var matches []ckanResource
	for _, res := range payload.Result.Resources {
		if strings.ToUpper(strings.TrimSpace(res.Format)) != wantFormat {
			continue
		}
		if strings.TrimSpace(res.URL) == "" {
			continue
		}
		matches = append(matches, res)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no %s resource in ckan package %s", wantFormat, packageID)
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].LastModified > matches[j].LastModified
	})

	return matches[0].URL, nil
}
