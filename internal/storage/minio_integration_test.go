package storage_test

import (
	"context"
	"os"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/storage"
)

func TestMinIOListPrefixIntegration(t *testing.T) {
	if os.Getenv("MINIO_INTEGRATION") != "1" {
		t.Skip("set MINIO_INTEGRATION=1 with docker compose minio running")
	}

	cfg := config.Config{
		StorageMode:    config.StorageModeMinIO,
		MinIOEndpoint:  envOr("MINIO_ENDPOINT", "http://localhost:9000"),
		MinIOAccessKey: envOr("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: envOr("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:    envOr("MINIO_BUCKET", "open-data-agro"),
	}

	store, err := storage.NewBronzeStore(cfg)
	if err != nil {
		t.Fatalf("NewBronzeStore: %v", err)
	}

	ctx := context.Background()
	prefix := "bronze/integration-test/"
	key := prefix + "part-smoke.parquet"
	payload := []byte("PAR1integration")

	if err := store.Put(ctx, key, payload, "application/vnd.apache.parquet"); err != nil {
		t.Fatalf("Put: %v", err)
	}
	t.Cleanup(func() { _ = store.Delete(ctx, key) })

	keys, err := store.ListPrefix(ctx, prefix)
	if err != nil {
		t.Fatalf("ListPrefix: %v", err)
	}
	found := false
	for _, listed := range keys {
		if listed == key {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("key %q not in listing: %v", key, keys)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
