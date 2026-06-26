package storage_test

import (
	"context"
	"os"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/storage"
)

func TestR2ListPrefixIntegration(t *testing.T) {
	if os.Getenv("R2_INTEGRATION") != "1" {
		t.Skip("set R2_INTEGRATION=1 with production R2 credentials in .env")
	}

	accountID := envOr("R2_ACCOUNT_ID", "")
	endpoint := envOr("R2_ENDPOINT", "")
	if endpoint == "" && accountID != "" {
		endpoint = "https://" + accountID + ".r2.cloudflarestorage.com"
	}

	cfg := config.Config{
		StorageMode:       config.StorageModeR2,
		R2AccessKeyID:     envOr("R2_ACCESS_KEY_ID", ""),
		R2SecretAccessKey: envOr("R2_SECRET_ACCESS_KEY", ""),
		R2Bucket:          envOr("R2_BUCKET", "open-data-agro"),
		R2Endpoint:        endpoint,
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
