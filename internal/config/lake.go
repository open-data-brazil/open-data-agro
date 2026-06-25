package config

import (
	"fmt"
	"os"
	"strings"
)

// LakeConfig holds lake path settings for the processor (no database required).
type LakeConfig struct {
	StorageMode       string
	LakeLocalRoot     string
	DeltaStoragePath  string
	DuckDBPath        string
	DeltaMinVersions  int
	MinIOEndpoint     string
	MinIOAccessKey    string
	MinIOSecretKey    string
	MinIOBucket       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2Bucket          string
	R2Endpoint        string
}

// LoadLakeFromEnv reads lake and storage settings for cmd/processor.
func LoadLakeFromEnv() (LakeConfig, error) {
	cfg := LakeConfig{
		StorageMode:       strings.TrimSpace(os.Getenv("STORAGE_MODE")),
		LakeLocalRoot:     strings.TrimSpace(os.Getenv("LAKE_LOCAL_ROOT")),
		DeltaStoragePath:  strings.TrimSpace(os.Getenv("DELTA_STORAGE_PATH")),
		DuckDBPath:        strings.TrimSpace(os.Getenv("DUCKDB_PATH")),
		MinIOEndpoint:     strings.TrimSpace(os.Getenv("MINIO_ENDPOINT")),
		MinIOAccessKey:    strings.TrimSpace(os.Getenv("MINIO_ACCESS_KEY")),
		MinIOSecretKey:    strings.TrimSpace(os.Getenv("MINIO_SECRET_KEY")),
		MinIOBucket:       strings.TrimSpace(os.Getenv("MINIO_BUCKET")),
		R2AccessKeyID:     strings.TrimSpace(os.Getenv("R2_ACCESS_KEY_ID")),
		R2SecretAccessKey: strings.TrimSpace(os.Getenv("R2_SECRET_ACCESS_KEY")),
		R2Bucket:          strings.TrimSpace(os.Getenv("R2_BUCKET")),
		R2Endpoint:        strings.TrimSpace(os.Getenv("R2_ENDPOINT")),
	}

	if cfg.StorageMode == "" {
		cfg.StorageMode = StorageModeLocal
	}
	cfg.StorageMode = strings.ToLower(cfg.StorageMode)

	if cfg.LakeLocalRoot == "" {
		cfg.LakeLocalRoot = "./lake"
	}
	if cfg.DeltaStoragePath == "" {
		cfg.DeltaStoragePath = "./lake/silver/"
	}
	if cfg.DuckDBPath == "" {
		cfg.DuckDBPath = "./duckdb/analytics.duckdb"
	}

	cfg.DeltaMinVersions = 30
	if raw := strings.TrimSpace(os.Getenv("DELTA_MIN_VERSIONS")); raw != "" {
		var n int
		if _, err := fmt.Sscanf(raw, "%d", &n); err != nil || n < 1 {
			return cfg, fmt.Errorf("DELTA_MIN_VERSIONS must be a positive integer")
		}
		cfg.DeltaMinVersions = n
	}

	accountID := strings.TrimSpace(os.Getenv("R2_ACCOUNT_ID"))
	cfg.R2Endpoint = resolveR2Endpoint(cfg.R2Endpoint, accountID)

	storageCfg := Config{
		StorageMode:       cfg.StorageMode,
		LakeLocalRoot:     cfg.LakeLocalRoot,
		MinIOEndpoint:     cfg.MinIOEndpoint,
		MinIOAccessKey:    cfg.MinIOAccessKey,
		MinIOSecretKey:    cfg.MinIOSecretKey,
		MinIOBucket:       cfg.MinIOBucket,
		R2AccessKeyID:     cfg.R2AccessKeyID,
		R2SecretAccessKey: cfg.R2SecretAccessKey,
		R2Bucket:          cfg.R2Bucket,
		R2Endpoint:        cfg.R2Endpoint,
	}
	if err := storageCfg.validateStorage(); err != nil {
		return cfg, err
	}

	return cfg, nil
}
