package catalog

import (
	"fmt"
	"strings"
)

// SplitDatasetID returns agency prefix and slug (e.g. conab, estimativa-graos).
func SplitDatasetID(datasetID string) (agency, slug string, err error) {
	dot := strings.Index(datasetID, ".")
	if dot <= 0 || dot >= len(datasetID)-1 {
		return "", "", fmt.Errorf("invalid dataset id %q", datasetID)
	}
	return datasetID[:dot], datasetID[dot+1:], nil
}

// Agency returns the catalog agency segment of a dataset ID.
func Agency(datasetID string) (string, error) {
	agency, _, err := SplitDatasetID(datasetID)
	return agency, err
}

// StorageAgency returns the lake path agency segment (may differ from catalog agency).
func StorageAgency(datasetID string) (string, error) {
	agency, _, err := SplitDatasetID(datasetID)
	if err != nil {
		return "", err
	}
	switch agency {
	case "oecd-fao":
		return "oecd", nil
	default:
		return agency, nil
	}
}
