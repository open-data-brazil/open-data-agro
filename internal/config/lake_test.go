package config

import (
	"testing"
)

func TestLoadLakeFromEnvDefaults(t *testing.T) {
	t.Setenv("STORAGE_MODE", "")
	t.Setenv("LAKE_LOCAL_ROOT", "")
	t.Setenv("DELTA_STORAGE_PATH", "")

	cfg, err := LoadLakeFromEnv()
	if err != nil {
		t.Fatalf("LoadLakeFromEnv: %v", err)
	}
	if cfg.StorageMode != StorageModeLocal {
		t.Fatalf("StorageMode: got %q", cfg.StorageMode)
	}
	if cfg.LakeLocalRoot != "./lake" {
		t.Fatalf("LakeLocalRoot: got %q", cfg.LakeLocalRoot)
	}
	if cfg.DeltaStoragePath != "./lake/silver/" {
		t.Fatalf("DeltaStoragePath: got %q", cfg.DeltaStoragePath)
	}
	if cfg.DeltaMinVersions != 30 {
		t.Fatalf("DeltaMinVersions: got %d", cfg.DeltaMinVersions)
	}
}
