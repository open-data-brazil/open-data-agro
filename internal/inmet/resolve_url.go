package inmet

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const annualZIPBase = "https://portal.inmet.gov.br/uploads/dadoshistoricos"

// ResolveURL validates the catalog source URL for an INMET dataset.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return "", fmt.Errorf("dataset %s has no source_url", entry.DatasetID)
	}

	if strings.Contains(raw, "{year}") {
		return raw, nil
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid source_url for %s: %w", entry.DatasetID, err)
	}
	if parsed.Scheme != "https" {
		return "", fmt.Errorf("source_url for %s must use https", entry.DatasetID)
	}
	host := strings.ToLower(parsed.Host)
	if host != "portal.inmet.gov.br" && host != "bdmep.inmet.gov.br" && host != "apitempo.inmet.gov.br" {
		return "", fmt.Errorf("source_url for %s must be on portal.inmet.gov.br, bdmep.inmet.gov.br, or apitempo.inmet.gov.br", entry.DatasetID)
	}
	return parsed.String(), nil
}

// AnnualZIPURL returns the BDMEP annual bulk ZIP URL for automatic stations.
func AnnualZIPURL(year int) string {
	return fmt.Sprintf("%s/%d.zip", annualZIPBase, year)
}

// ResolveAnnualZIPURL substitutes {year} in a catalog template or builds the default URL.
func ResolveAnnualZIPURL(entry catalog.RegistryEntry, year int) (string, error) {
	if year <= 0 {
		return "", fmt.Errorf("year is required for %s", entry.DatasetID)
	}
	if year < 2000 {
		return "", fmt.Errorf("INMET annual data starts at 2000, got %d", year)
	}

	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return AnnualZIPURL(year), nil
	}
	if strings.Contains(raw, "{year}") {
		return strings.ReplaceAll(raw, "{year}", strconv.Itoa(year)), nil
	}
	return raw, nil
}
