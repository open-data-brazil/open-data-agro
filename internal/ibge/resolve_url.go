package ibge

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const localidadesDocsURL = "https://servicodados.ibge.gov.br/api/docs/localidades"

// LocalidadesDocsURL is the official IBGE Localidades API documentation URL.
const LocalidadesDocsURL = localidadesDocsURL

// ResolveURL validates and returns the catalog source URL for an IBGE dataset.
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
	if host != "servicodados.ibge.gov.br" {
		return "", fmt.Errorf("source_url for %s must be on servicodados.ibge.gov.br", entry.DatasetID)
	}

	return parsed.String(), nil
}
