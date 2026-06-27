package mapa

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

const defaultCKANPackageShowURL = "https://dados.agricultura.gov.br/api/3/action/package_show"

// ckanPackageShowURL is the CKAN package_show endpoint (overridable in tests).
var ckanPackageShowURL = defaultCKANPackageShowURL

var safraYearPattern = regexp.MustCompile(`Safra (\d{4})/(\d{4})`)

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

// ListCKANResources returns CKAN resources for a MAPA package id.
func ListCKANResources(ctx context.Context, packageID string) ([]ckanResource, error) {
	return fetchCKANResources(ctx, packageID)
}

// ResolveURL returns the CKAN resource download URL for a MAPA catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	id := entry.DatasetID.String()
	if id == "mapa.sif-abate-estatisticas" {
		return ResolveSIFAbateURL(entry)
	}
	if strings.HasPrefix(id, "mapa.sipeagro-") {
		return resolveSIPEAGROURL(entry)
	}
	if id == "mapa.sisser-seguro-rural" {
		return resolveSISSERURL(entry)
	}

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

	nameFilter := strings.TrimSpace(entry.CKANResourceNameContains)
	if nameFilter != "" {
		return resolveCKANResourceByName(ctx, packageID, format, nameFilter)
	}

	return resolveLatestSafraResourceURL(ctx, packageID, format)
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

func resolveLatestSafraResourceURL(ctx context.Context, packageID, format string) (string, error) {
	resources, err := fetchCKANResources(ctx, packageID)
	if err != nil {
		return "", err
	}

	wantFormat := strings.ToUpper(strings.TrimSpace(format))
	type candidate struct {
		resource ckanResource
		safraKey int
	}
	var matches []candidate
	for _, res := range resources {
		if strings.ToUpper(strings.TrimSpace(res.Format)) != wantFormat {
			continue
		}
		if strings.TrimSpace(res.URL) == "" {
			continue
		}
		key := safraSortKey(res.Name)
		if key < 0 {
			continue
		}
		matches = append(matches, candidate{resource: res, safraKey: key})
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no annual safra %s resource in ckan package %s", wantFormat, packageID)
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].safraKey != matches[j].safraKey {
			return matches[i].safraKey > matches[j].safraKey
		}
		return matches[i].resource.LastModified > matches[j].resource.LastModified
	})

	return matches[0].resource.URL, nil
}

func resolveCKANResourceByName(ctx context.Context, packageID, format, nameContains string) (string, error) {
	resources, err := fetchCKANResources(ctx, packageID)
	if err != nil {
		return "", err
	}

	wantFormat := strings.ToUpper(strings.TrimSpace(format))
	needle := strings.ToLower(strings.TrimSpace(nameContains))
	var matches []ckanResource
	for _, res := range resources {
		if strings.ToUpper(strings.TrimSpace(res.Format)) != wantFormat {
			continue
		}
		if strings.TrimSpace(res.URL) == "" {
			continue
		}
		if !strings.Contains(strings.ToLower(res.Name), needle) {
			continue
		}
		matches = append(matches, res)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no %s resource matching %q in ckan package %s", wantFormat, nameContains, packageID)
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].LastModified > matches[j].LastModified
	})

	return matches[0].URL, nil
}

func fetchCKANResources(ctx context.Context, packageID string) ([]ckanResource, error) {
	url := fmt.Sprintf("%s?id=%s", ckanPackageShowURL, packageID)

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

func safraSortKey(name string) int {
	match := safraYearPattern.FindStringSubmatch(name)
	if len(match) < 2 {
		return -1
	}
	year, err := strconv.Atoi(match[1])
	if err != nil {
		return -1
	}
	return year
}
