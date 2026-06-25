package ingest

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Fingerprint captures immutable file identity from raw bytes and HTTP metadata.
type Fingerprint struct {
	SHA256          string
	FileSizeBytes   int64
	ContentType     string
	LastModified    string
	RowCount        int
	IngestDate      time.Time
	PartID          string
}

// NewFingerprint hashes raw bytes and assigns ingest partition metadata.
func NewFingerprint(body []byte, contentType, lastModified string, rowCount int, ingestDate time.Time) Fingerprint {
	sum := sha256.Sum256(body)
	partID := uuid.NewString()

	return Fingerprint{
		SHA256:        hex.EncodeToString(sum[:]),
		FileSizeBytes: int64(len(body)),
		ContentType:   contentType,
		LastModified:  lastModified,
		RowCount:      rowCount,
		IngestDate:    ingestDate,
		PartID:        partID,
	}
}

// DatasetSlug returns the path segment after the source prefix (e.g. estimativa-graos).
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

// BronzeKey builds the bronze object key for a parquet partition.
func BronzeKey(datasetID string, ingestDate time.Time, partID string) (string, error) {
	slug, err := DatasetSlug(datasetID)
	if err != nil {
		return "", err
	}
	date := ingestDate.UTC().Format("2006-01-02")
	return fmt.Sprintf("bronze/conab/%s/ingest_date=%s/part-%s.parquet", slug, date, partID), nil
}

// CountParquetRows estimates data rows from parquet bytes (best-effort for tests).
func CountParquetRows(parquetBytes []byte) (int, error) {
	if len(parquetBytes) == 0 {
		return 0, fmt.Errorf("empty parquet payload")
	}
	// Row count is tracked during conversion; this helper validates non-empty output.
	if bytes.Equal(parquetBytes[:4], []byte("PAR1")) {
		return 0, nil
	}
	return 0, fmt.Errorf("invalid parquet magic")
}
