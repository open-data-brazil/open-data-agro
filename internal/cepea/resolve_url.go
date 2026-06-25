package cepea

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL validates the catalog portal URL for a CEPEA dataset.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourcePortalURL)
	if raw == "" {
		var err error
		raw, err = PortalIndicatorURL(entry.DatasetID.String())
		if err != nil {
			return "", err
		}
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid source_portal_url for %s: %w", entry.DatasetID, err)
	}
	if parsed.Scheme != "https" {
		return "", fmt.Errorf("source_portal_url for %s must use https", entry.DatasetID)
	}
	host := strings.ToLower(parsed.Host)
	if !strings.Contains(host, "cepea.org.br") && !strings.Contains(host, "cepea.esalq.usp.br") {
		return "", fmt.Errorf("source_portal_url for %s must be on cepea.org.br", entry.DatasetID)
	}
	return parsed.String(), nil
}
