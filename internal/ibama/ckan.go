package ibama

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultCKANPackageShowURL = "https://dadosabertos.ibama.gov.br/api/3/action/package_show"

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

func fetchCKANResources(ctx context.Context, packageID string) ([]ckanResource, error) {
	url := fmt.Sprintf("%s?id=%s", defaultCKANPackageShowURL, packageID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ckan package_show: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ckan status %d for package %s", resp.StatusCode, packageID)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}

	var payload ckanPackageResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("parse ckan response: %w", err)
	}
	if !payload.Success {
		return nil, fmt.Errorf("ckan package_show failed for %s", packageID)
	}
	return payload.Result.Resources, nil
}

func resolveCKANResourceURL(ctx context.Context, packageID, resourceName, format string) (string, error) {
	resources, err := fetchCKANResources(ctx, packageID)
	if err != nil {
		return "", err
	}

	wantFormat := strings.ToUpper(strings.TrimSpace(format))
	needle := strings.ToLower(strings.TrimSpace(resourceName))
	for _, res := range resources {
		if strings.ToUpper(strings.TrimSpace(res.Format)) != wantFormat {
			continue
		}
		if strings.TrimSpace(res.URL) == "" {
			continue
		}
		if needle != "" && !strings.EqualFold(strings.TrimSpace(res.Name), resourceName) {
			continue
		}
		return res.URL, nil
	}
	return "", fmt.Errorf("no %s resource %q in ckan package %s", wantFormat, resourceName, packageID)
}

func resolveDirectURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return "", fmt.Errorf("dataset %s has no source_url", entry.DatasetID)
	}
	lower := strings.ToLower(raw)
	if !strings.Contains(lower, "gov.br") && !strings.Contains(lower, "blob.core.windows.net") {
		return "", fmt.Errorf("source_url for %s must be on gov.br or Azure blob", entry.DatasetID)
	}
	return raw, nil
}
