package ingest_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/ingest"
	"github.com/open-data-brazil/open-data-agro/internal/storage"
)

func TestBronzeLakeLocalWritePath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if prev := os.Getenv("LAKE_LOCAL_ROOT"); prev != "" {
		t.Cleanup(func() { _ = os.Setenv("LAKE_LOCAL_ROOT", prev) })
	}
	_ = os.Setenv("LAKE_LOCAL_ROOT", root)

	store := storage.NewLocalBronzeStoreForTest(root)
	ctx := context.Background()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}
	raw := []byte("col_a;col_b\n1;foo\n2;bar\n")

	parquetBytes, rowCount, err := ingest.ConvertToParquet(entry, raw)
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	if rowCount != 2 {
		t.Fatalf("rowCount: got %d want 2", rowCount)
	}

	ingestDate := time.Date(2026, 6, 25, 12, 0, 0, 0, time.UTC)
	partID := "test-part"
	bronzeKey, err := ingest.BronzeKey(entry.DatasetID.String(), ingestDate, partID)
	if err != nil {
		t.Fatalf("BronzeKey: %v", err)
	}
	if err := store.Put(ctx, bronzeKey, parquetBytes, "application/vnd.apache.parquet"); err != nil {
		t.Fatalf("Put parquet: %v", err)
	}

	fp := ingest.Fingerprint{
		SHA256:     "deadbeef",
		RowCount:   rowCount,
		IngestDate: ingestDate,
		PartID:     partID,
	}
	meta := ingest.NewBronzeMetadata(entry, fp, "https://portaldeinformacoes.conab.gov.br/download-arquivos.html", store.Backend())
	metaBytes, err := ingest.MarshalBronzeMetadata(meta)
	if err != nil {
		t.Fatalf("MarshalBronzeMetadata: %v", err)
	}
	metadataKey, err := ingest.PartitionMetadataKey(entry.DatasetID.String(), ingestDate)
	if err != nil {
		t.Fatalf("PartitionMetadataKey: %v", err)
	}
	if err := store.Put(ctx, metadataKey, metaBytes, "application/json"); err != nil {
		t.Fatalf("Put metadata: %v", err)
	}

	parquetPath := filepath.Join(root, filepath.FromSlash(bronzeKey))
	if _, err := os.Stat(parquetPath); err != nil {
		t.Fatalf("parquet missing: %v", err)
	}
	metaPath := filepath.Join(root, filepath.FromSlash(metadataKey))
	metaOnDisk, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("read metadata: %v", err)
	}
	var decoded ingest.BronzeMetadata
	if err := json.Unmarshal(metaOnDisk, &decoded); err != nil {
		t.Fatalf("decode metadata: %v", err)
	}
	if decoded.RowCount != 2 || decoded.StorageMode != "local" {
		t.Fatalf("unexpected metadata: %+v", decoded)
	}

	prefix := "bronze/conab/estimativa-graos/ingest_date=2026-06-25/"
	keys, err := store.ListPrefix(ctx, prefix)
	if err != nil {
		t.Fatalf("ListPrefix: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("ListPrefix: got %v", keys)
	}
	if !strings.HasPrefix(string(parquetBytes), "PAR1") {
		t.Fatal("expected parquet magic bytes")
	}
}
