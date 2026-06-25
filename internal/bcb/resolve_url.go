package bcb

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const sgsAPIBase = "https://api.bcb.gov.br/dados/serie/bcdata.sgs"

// ResolveURL validates the catalog base URL for a BCB SGS dataset.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return "", fmt.Errorf("dataset %s has no source_url", entry.DatasetID)
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid source_url for %s: %w", entry.DatasetID, err)
	}
	if parsed.Scheme != "https" {
		return "", fmt.Errorf("source_url for %s must use https", entry.DatasetID)
	}
	if !strings.EqualFold(parsed.Host, "api.bcb.gov.br") {
		return "", fmt.Errorf("source_url for %s must be on api.bcb.gov.br", entry.DatasetID)
	}
	return parsed.String(), nil
}

// SeriesBaseURL returns the SGS series endpoint without query parameters.
func SeriesBaseURL(code int) string {
	return fmt.Sprintf("%s.%d/dados", sgsAPIBase, code)
}
