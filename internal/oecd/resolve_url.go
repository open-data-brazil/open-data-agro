package oecd

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultSDMXBase = "https://sdmx.oecd.org/public/rest/data"

// ResolveURL returns the OECD SDMX data URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	url, err := buildCommodityURL(entry, "")
	if err != nil {
		return "", err
	}
	return url + "?format=csvfilewithlabels", nil
}

func buildDataflowPath(entry catalog.RegistryEntry) (string, error) {
	agency := strings.TrimSpace(entry.OECDAgency)
	datasetID := strings.TrimSpace(entry.OECDDatasetID)
	version := strings.TrimSpace(entry.OECDDatasetVersion)
	if agency == "" || datasetID == "" || version == "" {
		return "", fmt.Errorf("dataset %s missing oecd_agency, oecd_dataset_id, or oecd_dataset_version", entry.DatasetID)
	}
	return fmt.Sprintf("%s/%s,%s,%s", defaultSDMXBase, agency, datasetID, version), nil
}

func buildCommodityURL(entry catalog.RegistryEntry, commodity string) (string, error) {
	base, err := buildDataflowPath(entry)
	if err != nil {
		return "", err
	}
	selection := strings.TrimSpace(entry.OECDDataSelection)
	if selection == "" {
		selection = "BRA"
	}
	versionID := outlookVersionID(entry)
	commodityKey := strings.TrimSpace(commodity)
	if commodityKey == "" {
		commodityKey = "."
	}
	return fmt.Sprintf("%s/%s.A.%s..T.%s", base, selection, commodityKey, versionID), nil
}

func outlookVersionID(entry catalog.RegistryEntry) string {
	if raw := strings.TrimSpace(entry.SourceURL); strings.Contains(raw, "AO_") {
		if idx := strings.LastIndex(raw, "AO_"); idx >= 0 {
			segment := raw[idx:]
			if end := strings.IndexAny(segment, "/?"); end > 0 {
				return segment[:end]
			}
			return segment
		}
	}
	return "AO_2024_2033"
}
