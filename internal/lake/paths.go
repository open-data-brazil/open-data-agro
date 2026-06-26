// Package lake defines local-first path conventions for bronze, silver, and gold layers.
package lake

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// DatasetSlug returns the catalog slug segment (e.g. estimativa-graos).
func DatasetSlug(datasetID string) (string, error) {
	_, slug, err := catalog.SplitDatasetID(datasetID)
	return slug, err
}

// DatasetAgency returns the catalog agency segment (e.g. conab, anp).
func DatasetAgency(datasetID string) (string, error) {
	agency, _, err := catalog.SplitDatasetID(datasetID)
	return agency, err
}

// SilverTableName maps a dataset ID to the Delta table directory name (underscores).
func SilverTableName(datasetID string) (string, error) {
	slug, err := DatasetSlug(datasetID)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(slug, "-", "_"), nil
}

// BronzeGlob returns a glob for all bronze parquet files of a dataset.
func BronzeGlob(lakeRoot, datasetID string) (string, error) {
	agency, slug, err := storageParts(datasetID)
	if err != nil {
		return "", err
	}
	return filepath.Join(lakeRoot, "bronze", agency, slug, "**", "*.parquet"), nil
}

// BronzeDir returns the bronze dataset directory.
func BronzeDir(lakeRoot, datasetID string) (string, error) {
	agency, slug, err := storageParts(datasetID)
	if err != nil {
		return "", err
	}
	return filepath.Join(lakeRoot, "bronze", agency, slug), nil
}

// SilverTableDir returns the local Delta silver table directory.
func SilverTableDir(silverRoot, datasetID string) (string, error) {
	agency, table, err := silverParts(datasetID)
	if err != nil {
		return "", err
	}
	return filepath.Join(silverRoot, agency, table), nil
}

// GoldMartDir returns the dbt gold mart directory for a dataset.
func GoldMartDir(goldRoot, datasetID string) (string, error) {
	agency, table, err := silverParts(datasetID)
	if err != nil {
		return "", err
	}
	return filepath.Join(goldRoot, fmt.Sprintf("mart_%s__%s", agency, table)), nil
}

func silverParts(datasetID string) (agency, table string, err error) {
	agency, slug, err := storageParts(datasetID)
	if err != nil {
		return "", "", err
	}
	table = strings.ReplaceAll(slug, "-", "_")
	return agency, table, nil
}

func storageParts(datasetID string) (agency, slug string, err error) {
	agency, err = catalog.StorageAgency(datasetID)
	if err != nil {
		return "", "", err
	}
	_, slug, err = catalog.SplitDatasetID(datasetID)
	return agency, slug, err
}

// NormalizeRoot trims trailing slashes from a lake layer root.
func NormalizeRoot(root string) string {
	return strings.TrimRight(strings.TrimSpace(root), "/")
}
