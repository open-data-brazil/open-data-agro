package processor

import (
	"context"
	"os"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
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
	reg := catalog.NewRegistry([]catalog.RegistryEntry{{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
	}})

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

	_, err = smoker.Smoke(context.Background(), SmokeOptions{DatasetID: "conab.estimativa-graos"})
	if err != nil {
		t.Fatalf("Smoke against MinIO: %v", err)
	}
}
