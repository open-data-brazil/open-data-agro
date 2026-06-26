package processor

import (
	"context"
	"os"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/ingest"
	"github.com/open-data-brazil/open-data-agro/internal/storage"
)

func TestSmokeMinIOPathSwap(t *testing.T) {
	if os.Getenv("MINIO_INTEGRATION") != "1" {
		t.Skip("set MINIO_INTEGRATION=1 with docker compose minio running")
	}

	cfg := config.LakeConfig{
		StorageMode:    config.StorageModeMinIO,
		LakeLocalRoot:  "./lake",
		DuckDBPath:     ":memory:",
		MinIOEndpoint:  envOr("MINIO_ENDPOINT", "http://localhost:9000"),
		MinIOAccessKey: envOr("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: envOr("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:    envOr("MINIO_BUCKET", "open-data-agro"),
	}
	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	}
	reg := catalog.NewRegistry([]catalog.RegistryEntry{entry})

	smoker, err := NewSmoker(cfg, reg)
	if err != nil {
		t.Fatalf("NewSmoker: %v", err)
	}

	uri, err := smoker.bronzeURI(SmokeOptions{DatasetID: "conab.estimativa-graos"})
	if err != nil {
		t.Fatalf("bronzeURI: %v", err)
	}
	wantPrefix := "s3://open-data-agro/bronze/conab/estimativa-graos/"
	if uri != wantPrefix+"**/*.parquet" {
		t.Fatalf("got %q want s3 bronze glob", uri)
	}

	parquetBytes, _, err := ingest.ConvertToParquet(entry, []byte("a;b\n1;x\n2;y\n"))
	if err != nil {
		t.Fatalf("ConvertToParquet: %v", err)
	}
	store, err := storage.NewBronzeStore(config.Config{
		StorageMode:    cfg.StorageMode,
		MinIOEndpoint:  cfg.MinIOEndpoint,
		MinIOAccessKey: cfg.MinIOAccessKey,
		MinIOSecretKey: cfg.MinIOSecretKey,
		MinIOBucket:    cfg.MinIOBucket,
	})
	if err != nil {
		t.Fatalf("NewBronzeStore: %v", err)
	}
	ctx := context.Background()
	key := "bronze/conab/estimativa-graos/ingest_date=2026-06-25/part-ci.parquet"
	if err := store.Put(ctx, key, parquetBytes, "application/vnd.apache.parquet"); err != nil {
		t.Fatalf("Put bronze to MinIO: %v", err)
	}
	t.Cleanup(func() { _ = store.Delete(ctx, key) })

	result, err := smoker.Smoke(ctx, SmokeOptions{DatasetID: "conab.estimativa-graos"})
	if err != nil {
		t.Fatalf("Smoke against MinIO: %v", err)
	}
	if result.RowCount != 2 {
		t.Fatalf("row_count: got %d want 2", result.RowCount)
	}
}
