package anp

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL validates and returns the download URL for an ANP catalog entry.
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
	host := strings.ToLower(parsed.Host)
	if !strings.HasSuffix(host, "gov.br") {
		return "", fmt.Errorf("source_url for %s must be on gov.br host", entry.DatasetID)
	}

	if strings.HasSuffix(strings.ToLower(parsed.Path), ".xlsx") {
		return parsed.String(), nil
	}

	lpcFile := strings.TrimSpace(entry.ANPLPCFile)
	if lpcFile == "" {
		return "", fmt.Errorf("dataset %s requires anp_lpc_file or a direct .xlsx source_url", entry.DatasetID)
	}

	return ResolveLatestLPCURL(parsed.String(), lpcFile)
}
