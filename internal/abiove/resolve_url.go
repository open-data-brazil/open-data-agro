package abiove

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultPortalURL = "https://abiove.org.br/estatisticas/"

// ResolveURL validates the catalog download URL for an Abiove dataset.
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
	if !strings.Contains(host, "abiove.org.br") {
		return "", fmt.Errorf("source_url for %s must be on abiove.org.br", entry.DatasetID)
	}
	if !strings.HasSuffix(strings.ToLower(parsed.Path), ".xlsx") {
		return "", fmt.Errorf("source_url for %s must point to an .xlsx file", entry.DatasetID)
	}
	return parsed.String(), nil
}

// PortalURL returns the Abiove statistics landing page.
func PortalURL(entry catalog.RegistryEntry) string {
	raw := strings.TrimSpace(entry.SourcePortalURL)
	if raw != "" {
		return raw
	}
	return defaultPortalURL
}
