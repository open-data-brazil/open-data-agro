package oecd

import (
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// BuildMetadataProbeURL returns the SDMX dataflow metadata endpoint (lighter than full data pull).
func BuildMetadataProbeURL(entry catalog.RegistryEntry) (string, error) {
	agency := strings.TrimSpace(entry.OECDAgency)
	datasetID := strings.TrimSpace(entry.OECDDatasetID)
	version := strings.TrimSpace(entry.OECDDatasetVersion)
	if agency == "" || datasetID == "" || version == "" {
		return "", fmt.Errorf("dataset %s missing oecd_agency, oecd_dataset_id, or oecd_dataset_version", entry.DatasetID)
	}
	return fmt.Sprintf("https://sdmx.oecd.org/public/rest/dataflow/%s/%s/%s", agency, datasetID, version), nil
}
