package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds runtime settings loaded from environment variables.
type Config struct {
	DatabaseURL       string
	LakeLocalRoot     string
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2Bucket          string
	R2Endpoint        string
	AlertWebhookURL   string
}

// LoadFromEnv reads configuration from the process environment.
func LoadFromEnv() (Config, error) {
	cfg := Config{
		DatabaseURL:       strings.TrimSpace(os.Getenv("DATABASE_URL")),
		LakeLocalRoot:     strings.TrimSpace(os.Getenv("LAKE_LOCAL_ROOT")),
		R2AccountID:       strings.TrimSpace(os.Getenv("R2_ACCOUNT_ID")),
		R2AccessKeyID:     strings.TrimSpace(os.Getenv("R2_ACCESS_KEY_ID")),
		R2SecretAccessKey: strings.TrimSpace(os.Getenv("R2_SECRET_ACCESS_KEY")),
		R2Bucket:          strings.TrimSpace(os.Getenv("R2_BUCKET")),
		R2Endpoint:        strings.TrimSpace(os.Getenv("R2_ENDPOINT")),
		AlertWebhookURL:   strings.TrimSpace(os.Getenv("ALERT_WEBHOOK_URL")),
	}

	if cfg.LakeLocalRoot == "" {
		cfg.LakeLocalRoot = "./lake"
	}
	if cfg.DatabaseURL == "" {
		return cfg, fmt.Errorf("DATABASE_URL is required")
	}

	cfg.R2Endpoint = resolveR2Endpoint(cfg.R2Endpoint, cfg.R2AccountID)
	return cfg, nil
}

// R2Enabled reports whether R2 upload credentials are configured.
func (c Config) R2Enabled() bool {
	return c.R2AccessKeyID != "" && c.R2SecretAccessKey != "" && c.R2Bucket != "" && c.R2Endpoint != ""
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
