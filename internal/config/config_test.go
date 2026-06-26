package config

import (
	"os"
	"testing"
)

func TestLoadFromEnvDefaultsToLocalStorage(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://localhost/test")
	t.Setenv("STORAGE_MODE", "")
	t.Setenv("R2_ACCESS_KEY_ID", "should-not-select-r2")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv: %v", err)
	}
	if cfg.StorageMode != StorageModeLocal {
		t.Fatalf("StorageMode: got %q want %q", cfg.StorageMode, StorageModeLocal)
	}
}

func TestLoadFromEnvRequiresMinIOCredentials(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://localhost/test")
	t.Setenv("STORAGE_MODE", "minio")
	t.Setenv("MINIO_ENDPOINT", "")

	if _, err := LoadFromEnv(); err == nil {
		t.Fatal("expected error for incomplete minio config")
	}
}

func TestLoadFromEnvResolvesR2EndpointFromAccountID(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://localhost/test")
	t.Setenv("STORAGE_MODE", "r2")
	t.Setenv("R2_ACCOUNT_ID", "abc123")
	t.Setenv("R2_ACCESS_KEY_ID", "key")
	t.Setenv("R2_SECRET_ACCESS_KEY", "secret")
	t.Setenv("R2_BUCKET", "open-data-agro")
	t.Setenv("R2_ENDPOINT", "")

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv: %v", err)
	}
	want := "https://abc123.r2.cloudflarestorage.com"
	if cfg.R2Endpoint != want {
		t.Fatalf("R2Endpoint: got %q want %q", cfg.R2Endpoint, want)
	}
}

func TestLoadFromEnvRequiresR2Credentials(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://localhost/test")
	t.Setenv("STORAGE_MODE", "r2")
	t.Setenv("R2_ACCOUNT_ID", "")
	t.Setenv("R2_ACCESS_KEY_ID", "")
	t.Setenv("R2_SECRET_ACCESS_KEY", "secret")
	t.Setenv("R2_BUCKET", "open-data-agro")
	t.Setenv("R2_ENDPOINT", "")

	if _, err := LoadFromEnv(); err == nil {
		t.Fatal("expected error for incomplete r2 config")
	}
}

func TestValidateR2EnvLive(t *testing.T) {
	if os.Getenv("VALIDATE_R2_ENV") != "1" {
		t.Skip("set VALIDATE_R2_ENV=1 via scripts/deploy/validate_r2_env.sh")
	}

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv: %v", err)
	}
	if cfg.StorageMode != StorageModeR2 {
		t.Fatalf("StorageMode: got %q want %q", cfg.StorageMode, StorageModeR2)
	}
	if cfg.R2Bucket == "" || cfg.R2Endpoint == "" {
		t.Fatal("R2 bucket and endpoint must be set")
	}
	if cfg.R2AccessKeyID == "" || cfg.R2SecretAccessKey == "" {
		t.Fatal("R2 access credentials must be set")
	}
}
