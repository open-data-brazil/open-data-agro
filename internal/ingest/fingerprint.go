package ingest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
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
	partID, err := newUUIDv7String()
	if err != nil {
		panic(fmt.Sprintf("uuid v7: %v", err))
	}

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

// DatasetSlug returns the path segment after the agency prefix (e.g. estimativa-graos).
func DatasetSlug(datasetID string) (string, error) {
	_, slug, err := catalog.SplitDatasetID(datasetID)
	return slug, err
}

// BronzeKey builds the bronze object key for a parquet partition.
func BronzeKey(datasetID string, ingestDate time.Time, partID string) (string, error) {
	agency, slug, err := catalog.SplitDatasetID(datasetID)
	if err != nil {
		return "", err
	}
	date := ingestDate.UTC().Format("2006-01-02")
	return fmt.Sprintf("bronze/%s/%s/ingest_date=%s/part-%s.parquet", agency, slug, date, partID), nil
}

func newUUIDv7String() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
