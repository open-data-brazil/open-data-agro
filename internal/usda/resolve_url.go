package usda

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// ResolveURL returns a descriptive endpoint reference for audit metadata.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	if entry.DatasetID.String() == "usda.wasde" {
		raw := strings.TrimSpace(entry.SourceURL)
		if raw == "" {
			return wasdeESMISIndexURL, nil
		}
		return raw, nil
	}

	code := strings.TrimSpace(entry.PSDCommodityCode)
	if code == "" {
		return "", fmt.Errorf("dataset %s missing psd_commodity_code", entry.DatasetID)
	}
	return fmt.Sprintf("%s#getDatabyCommodityPerYear?commodity=%s", psdSOAPEndpoint, code), nil
}
