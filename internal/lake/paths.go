// Package lake defines local-first path conventions for bronze, silver, and gold layers.
package lake

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DatasetSlug returns the catalog slug segment (e.g. estimativa-graos).
func DatasetSlug(datasetID string) (string, error) {
	const prefix = "conab."
	if !strings.HasPrefix(datasetID, prefix) {
		return "", fmt.Errorf("unsupported dataset prefix in %s", datasetID)
	}
	slug := strings.TrimPrefix(datasetID, prefix)
	if slug == "" {
		return "", fmt.Errorf("empty dataset slug for %s", datasetID)
	}
	return slug, nil
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
	slug, err := DatasetSlug(datasetID)
	if err != nil {
		return "", err
	}
	return filepath.Join(lakeRoot, "bronze", "conab", slug, "**", "*.parquet"), nil
}

// BronzeDir returns the bronze dataset directory.
func BronzeDir(lakeRoot, datasetID string) (string, error) {
	slug, err := DatasetSlug(datasetID)
	if err != nil {
		return "", err
	}
	return filepath.Join(lakeRoot, "bronze", "conab", slug), nil
}

// SilverTableDir returns the local Delta silver table directory.
func SilverTableDir(silverRoot, datasetID string) (string, error) {
	table, err := SilverTableName(datasetID)
	if err != nil {
		return "", err
	}
	return filepath.Join(silverRoot, "conab", table), nil
}

// GoldMartDir returns the dbt gold mart directory for a dataset.
func GoldMartDir(goldRoot, datasetID string) (string, error) {
	table, err := SilverTableName(datasetID)
	if err != nil {
		return "", err
	}
	return filepath.Join(goldRoot, "mart_conab__"+table), nil
}

// NormalizeRoot trims trailing slashes from a lake layer root.
func NormalizeRoot(root string) string {
	return strings.TrimRight(strings.TrimSpace(root), "/")
}
