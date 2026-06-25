package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/config"
)

func TestLocalBronzeStorePutAndDelete(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	store := newLocalBronzeStore(root)
	ctx := context.Background()
	key := "bronze/conab/estimativa-graos/ingest_date=2026-06-25/part-1.parquet"
	payload := []byte("PAR1test")

	if err := store.Put(ctx, key, payload, "application/vnd.apache.parquet"); err != nil {
		t.Fatalf("Put: %v", err)
	}

	path := filepath.Join(root, filepath.FromSlash(key))
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != string(payload) {
		t.Fatalf("unexpected payload: %q", data)
	}

	if err := store.Delete(ctx, key); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected file removed, stat err=%v", err)
	}
}

func TestLocalBronzeStoreListPrefix(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	store := newLocalBronzeStore(root)
	ctx := context.Background()
	prefix := "bronze/conab/estimativa-graos/ingest_date=2026-06-25/"
	keys := []string{
		prefix + "part-a.parquet",
		prefix + "_metadata.json",
	}

	for _, key := range keys {
		if err := store.Put(ctx, key, []byte("x"), "application/octet-stream"); err != nil {
			t.Fatalf("Put %s: %v", key, err)
		}
	}

	listed, err := store.ListPrefix(ctx, prefix)
	if err != nil {
		t.Fatalf("ListPrefix: %v", err)
	}
	if len(listed) != len(keys) {
		t.Fatalf("listed %d keys, want %d: %v", len(listed), len(keys), listed)
	}
}

func TestNewBronzeStoreDefaultsToLocal(t *testing.T) {
	t.Parallel()

	store, err := NewBronzeStore(config.Config{
		StorageMode:   config.StorageModeLocal,
		LakeLocalRoot: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("NewBronzeStore: %v", err)
	}
	if store.Backend() != config.StorageModeLocal {
		t.Fatalf("backend: got %q want local", store.Backend())
	}
}

func TestNewBronzeStoreSelectsMinIO(t *testing.T) {
	t.Parallel()

	store, err := NewBronzeStore(config.Config{
		StorageMode:    config.StorageModeMinIO,
		MinIOEndpoint:  "http://localhost:9000",
		MinIOAccessKey: "minioadmin",
		MinIOSecretKey: "minioadmin",
		MinIOBucket:    "open-data-agro",
	})
	if err != nil {
		t.Fatalf("NewBronzeStore: %v", err)
	}
	if store.Backend() != config.StorageModeMinIO {
		t.Fatalf("backend: got %q want minio", store.Backend())
	}
}
