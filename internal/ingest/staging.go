package ingest

import (
	"fmt"
	"os"
	"path/filepath"
)

// StagedDownload is a raw portal file written to a temp directory before conversion.
type StagedDownload struct {
	Path    string
	Cleanup func()
}

// StageRaw writes downloaded bytes to a private temp file for conversion.
func StageRaw(jobID string, body []byte) (*StagedDownload, error) {
	dir := filepath.Join(os.TempDir(), "open-data-agro-ingestor")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("create staging dir: %w", err)
	}

	path := filepath.Join(dir, fmt.Sprintf("%s.raw", jobID))
	if err := os.WriteFile(path, body, 0o600); err != nil {
		return nil, fmt.Errorf("write staged download: %w", err)
	}

	return &StagedDownload{
		Path: path,
		Cleanup: func() {
			_ = os.Remove(path)
		},
	}, nil
}

// ReadStaged loads staged bytes from disk.
func ReadStaged(staged *StagedDownload) ([]byte, error) {
	data, err := os.ReadFile(staged.Path)
	if err != nil {
		return nil, fmt.Errorf("read staged download: %w", err)
	}
	return data, nil
}
