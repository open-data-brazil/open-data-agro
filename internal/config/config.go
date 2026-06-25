package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	StorageModeLocal = "local"
	StorageModeMinIO = "minio"
	StorageModeR2    = "r2"
)

// Config holds runtime settings loaded from environment variables.
type Config struct {
	DatabaseURL       string
	StorageMode         string
	LakeLocalRoot       string
	MinIOEndpoint       string
	MinIOAccessKey      string
	MinIOSecretKey      string
	MinIOBucket         string
	R2AccountID         string
	R2AccessKeyID       string
	R2SecretAccessKey   string
	R2Bucket            string
	R2Endpoint          string
	AlertWebhookURL     string
}

// LoadFromEnv reads configuration from the process environment.
func LoadFromEnv() (Config, error) {
	cfg := Config{
		DatabaseURL:       strings.TrimSpace(os.Getenv("DATABASE_URL")),
		StorageMode:         strings.TrimSpace(os.Getenv("STORAGE_MODE")),
		LakeLocalRoot:       strings.TrimSpace(os.Getenv("LAKE_LOCAL_ROOT")),
		MinIOEndpoint:       strings.TrimSpace(os.Getenv("MINIO_ENDPOINT")),
		MinIOAccessKey:      strings.TrimSpace(os.Getenv("MINIO_ACCESS_KEY")),
		MinIOSecretKey:      strings.TrimSpace(os.Getenv("MINIO_SECRET_KEY")),
		MinIOBucket:         strings.TrimSpace(os.Getenv("MINIO_BUCKET")),
		R2AccountID:         strings.TrimSpace(os.Getenv("R2_ACCOUNT_ID")),
		R2AccessKeyID:       strings.TrimSpace(os.Getenv("R2_ACCESS_KEY_ID")),
		R2SecretAccessKey:   strings.TrimSpace(os.Getenv("R2_SECRET_ACCESS_KEY")),
		R2Bucket:            strings.TrimSpace(os.Getenv("R2_BUCKET")),
		R2Endpoint:          strings.TrimSpace(os.Getenv("R2_ENDPOINT")),
		AlertWebhookURL:     strings.TrimSpace(os.Getenv("ALERT_WEBHOOK_URL")),
	}

	if cfg.StorageMode == "" {
		cfg.StorageMode = StorageModeLocal
	}
	cfg.StorageMode = strings.ToLower(cfg.StorageMode)

	if cfg.LakeLocalRoot == "" {
		cfg.LakeLocalRoot = "./lake"
	}
	if cfg.DatabaseURL == "" {
		return cfg, fmt.Errorf("DATABASE_URL is required")
	}

	cfg.R2Endpoint = resolveR2Endpoint(cfg.R2Endpoint, cfg.R2AccountID)

	if err := cfg.validateStorage(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c Config) validateStorage() error {
	switch c.StorageMode {
	case StorageModeLocal:
		return nil
	case StorageModeMinIO:
		if c.MinIOEndpoint == "" || c.MinIOAccessKey == "" || c.MinIOSecretKey == "" || c.MinIOBucket == "" {
			return fmt.Errorf("STORAGE_MODE=minio requires MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, and MINIO_BUCKET")
		}
		return nil
	case StorageModeR2:
		if c.R2AccessKeyID == "" || c.R2SecretAccessKey == "" || c.R2Bucket == "" || c.R2Endpoint == "" {
			return fmt.Errorf("STORAGE_MODE=r2 requires R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_BUCKET, and R2_ENDPOINT (or R2_ACCOUNT_ID)")
		}
		return nil
	default:
		return fmt.Errorf("invalid STORAGE_MODE %q (want local, minio, or r2)", c.StorageMode)
	}
}

func resolveR2Endpoint(endpoint, accountID string) string {
	if endpoint == "" {
		if accountID == "" {
			return ""
		}
		return fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)
	}
	return strings.ReplaceAll(endpoint, "{account_id}", accountID)
}
