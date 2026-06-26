package igc

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultGOIURL = "https://igc.int/_csv/igc__goi.xlsb"

// ResolveURL returns the official download URL for an IGC catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		switch entry.DatasetID.String() {
		case "igc.goi-index":
			raw = defaultGOIURL
		default:
			return "", fmt.Errorf("dataset %s has no source_url", entry.DatasetID)
		}
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid source_url for %s: %w", entry.DatasetID, err)
	}
	if !strings.EqualFold(parsed.Host, "www.igc.int") && !strings.EqualFold(parsed.Host, "igc.int") {
		return "", fmt.Errorf("source_url for %s must be on igc.int", entry.DatasetID)
	}

	return raw, nil
}
