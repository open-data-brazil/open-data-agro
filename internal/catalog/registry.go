package catalog

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type registryFile struct {
	Entries []RegistryEntry `yaml:"entries"`
}

// LoadRegistryDir loads and merges registry YAML files from a directory tree.
func LoadRegistryDir(dir string) ([]RegistryEntry, error) {
	var entries []RegistryEntry

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("read registry %s: %w", path, readErr)
		}

		var file registryFile
		if unmarshalErr := yaml.Unmarshal(data, &file); unmarshalErr != nil {
			return fmt.Errorf("parse registry %s: %w", path, unmarshalErr)
		}

		for i := range file.Entries {
			if _, parseErr := ParseDatasetID(file.Entries[i].DatasetID.String()); parseErr != nil {
				return fmt.Errorf("registry %s entry %d: %w", path, i, parseErr)
			}
		}

		entries = append(entries, file.Entries...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// Registry provides lookup helpers over loaded catalog entries.
type Registry struct {
	entries []RegistryEntry
	byID    map[DatasetID]RegistryEntry
}

// NewRegistry builds an in-memory registry from entries.
func NewRegistry(entries []RegistryEntry) *Registry {
	byID := make(map[DatasetID]RegistryEntry, len(entries))
	for _, entry := range entries {
		byID[entry.DatasetID] = entry
	}
	return &Registry{entries: entries, byID: byID}
}

// Entries returns a copy of all registry entries.
func (r *Registry) Entries() []RegistryEntry {
	out := make([]RegistryEntry, len(r.entries))
	copy(out, r.entries)
	return out
}

// Get returns an entry by dataset ID string, or false when unknown or invalid.
func (r *Registry) Get(datasetID string) (RegistryEntry, bool) {
	if !IsDatasetID(datasetID) {
		return RegistryEntry{}, false
	}
	entry, ok := r.byID[DatasetID(datasetID)]
	return entry, ok
}

// Require returns an entry or an error when the dataset is unknown.
func (r *Registry) Require(datasetID string) (RegistryEntry, error) {
	entry, ok := r.Get(datasetID)
	if !ok {
		return RegistryEntry{}, fmt.Errorf("unknown dataset: %s", datasetID)
	}
	return entry, nil
}

// ListIDs returns all dataset IDs in registration order.
func (r *Registry) ListIDs() []DatasetID {
	ids := make([]DatasetID, len(r.entries))
	for i, entry := range r.entries {
		ids[i] = entry.DatasetID
	}
	return ids
}
