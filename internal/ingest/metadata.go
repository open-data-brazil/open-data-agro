package ingest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// BronzeMetadata is the JSON sidecar written next to each bronze partition.
type BronzeMetadata struct {
	DatasetID      string `json:"dataset_id"`
	SHA256         string `json:"sha256"`
	SourceURL      string `json:"source_url"`
	IngestedAt     string `json:"ingested_at"`
	RowCount       int    `json:"row_count"`
	OriginalFormat string `json:"original_format"`
	StorageMode    string `json:"storage_mode"`
}

// PartitionMetadataKey builds the _metadata.json key for a bronze partition.
func PartitionMetadataKey(datasetID string, ingestDate time.Time) (string, error) {
	slug, err := DatasetSlug(datasetID)
	if err != nil {
		return "", err
	}
	date := ingestDate.UTC().Format("2006-01-02")
	return fmt.Sprintf("bronze/conab/%s/ingest_date=%s/_metadata.json", slug, date), nil
}

// NewBronzeMetadata builds the partition sidecar payload.
func NewBronzeMetadata(entry catalog.RegistryEntry, fp Fingerprint, sourceURL, storageMode string) BronzeMetadata {
	return BronzeMetadata{
		DatasetID:      entry.DatasetID.String(),
		SHA256:         fp.SHA256,
		SourceURL:      sourceURL,
		IngestedAt:     fp.IngestDate.UTC().Format(time.RFC3339),
		RowCount:       fp.RowCount,
		OriginalFormat: string(entry.Format),
		StorageMode:    storageMode,
	}
}

// MarshalBronzeMetadata serializes the sidecar JSON.
func MarshalBronzeMetadata(meta BronzeMetadata) ([]byte, error) {
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal bronze metadata: %w", err)
	}
	return append(data, '\n'), nil
}
