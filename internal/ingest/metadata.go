package ingest

import (
	"encoding/json"
	"fmt"
	"strings"
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
	Agencia        string `json:"agencia,omitempty"`
	FonteOficial   string `json:"fonte_oficial,omitempty"`
	FonteTipo      string `json:"fonte_tipo,omitempty"`
	Licenca        string `json:"licenca,omitempty"`
	PortalLabel    string `json:"portal_label,omitempty"`
}

// PartitionMetadataKey builds the _metadata.json key for a bronze partition.
func PartitionMetadataKey(datasetID string, ingestDate time.Time) (string, error) {
	agency, slug, err := catalog.SplitDatasetID(datasetID)
	if err != nil {
		return "", err
	}
	date := ingestDate.UTC().Format("2006-01-02")
	return fmt.Sprintf("bronze/%s/%s/ingest_date=%s/_metadata.json", agency, slug, date), nil
}

// NewBronzeMetadata builds the partition sidecar payload.
func NewBronzeMetadata(entry catalog.RegistryEntry, fp Fingerprint, sourceURL, storageMode string) BronzeMetadata {
	meta := BronzeMetadata{
		DatasetID:      entry.DatasetID.String(),
		SHA256:         fp.SHA256,
		SourceURL:      sourceURL,
		IngestedAt:     fp.IngestDate.UTC().Format(time.RFC3339),
		RowCount:       fp.RowCount,
		OriginalFormat: string(entry.Format),
		StorageMode:    storageMode,
		PortalLabel:    entry.PortalLabel,
	}
	if strings.HasPrefix(entry.DatasetID.String(), "cepea.") {
		meta.Agencia = "CEPEA"
		meta.FonteOficial = entry.PortalURL()
		meta.FonteTipo = entry.FonteTipo
		meta.Licenca = entry.License
		if meta.FonteTipo == "" {
			meta.FonteTipo = "referencia_mercado"
		}
		if meta.Licenca == "" {
			meta.Licenca = "CC BY-NC 4.0"
		}
	}
	return meta
}

// MarshalBronzeMetadata serializes the sidecar JSON.
func MarshalBronzeMetadata(meta BronzeMetadata) ([]byte, error) {
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal bronze metadata: %w", err)
	}
	return append(data, '\n'), nil
}
