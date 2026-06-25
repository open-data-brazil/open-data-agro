// Package catalog defines dataset identifiers, metadata, and registry loading.
package catalog

import (
	"fmt"
	"regexp"
)

// DatasetID is a branded dataset identifier (e.g. conab.estimativa-graos).
type DatasetID string

var datasetIDPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(\.[a-z][a-z0-9-]*)+$`)

// ParseDatasetID validates and returns a DatasetID.
func ParseDatasetID(value string) (DatasetID, error) {
	if !datasetIDPattern.MatchString(value) {
		return "", fmt.Errorf("invalid dataset ID: %s", value)
	}
	return DatasetID(value), nil
}

// MustParseDatasetID parses a DatasetID or panics.
func MustParseDatasetID(value string) DatasetID {
	id, err := ParseDatasetID(value)
	if err != nil {
		panic(err)
	}
	return id
}

// String returns the underlying string value.
func (id DatasetID) String() string {
	return string(id)
}

// IsDatasetID reports whether value matches the dataset ID pattern.
func IsDatasetID(value string) bool {
	return datasetIDPattern.MatchString(value)
}
