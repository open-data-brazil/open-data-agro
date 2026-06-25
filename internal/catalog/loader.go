package catalog

import (
	"os"
	"path/filepath"
)

// DefaultRegistryDir is the path to YAML catalog configs relative to the repo root.
const DefaultRegistryDir = "configs/catalog"

// LoadDefaultRegistry loads all YAML entries from configs/catalog.
func LoadDefaultRegistry() (*Registry, error) {
	root, err := repoRoot()
	if err != nil {
		return nil, err
	}
	return LoadRegistryFrom(filepath.Join(root, DefaultRegistryDir))
}

// LoadRegistryFrom loads registry YAML files from an absolute or relative directory.
func LoadRegistryFrom(dir string) (*Registry, error) {
	entries, err := LoadRegistryDir(dir)
	if err != nil {
		return nil, err
	}
	return NewRegistry(entries), nil
}

func repoRoot() (string, error) {
	if root := os.Getenv("OPEN_DATA_AGRO_ROOT"); root != "" {
		return root, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return wd, nil
		}
		dir = parent
	}
}
