package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
	"github.com/open-data-brazil/open-data-agro/internal/lake"
)

// PromoteOptions controls a bronze → silver promotion run.
type PromoteOptions struct {
	DatasetID string
}

// PromoteResult summarizes a promotion attempt.
type PromoteResult struct {
	DatasetID  string
	RowCount   int
	BronzeDir  string
	SilverDir  string
	StorageMode string
}

// Promoter runs bronze → silver Delta promotions.
type Promoter struct {
	cfg      config.LakeConfig
	registry *catalog.Registry
	python   string
	script   string
}

// NewPromoter wires promotion dependencies.
func NewPromoter(cfg config.LakeConfig, registry *catalog.Registry) *Promoter {
	return &Promoter{
		cfg:      cfg,
		registry: registry,
		python:   envOr("PYTHON", "python3"),
		script:   defaultPromoteScript(),
	}
}

// Promote reads bronze Parquet and appends a new Delta silver version.
func (p *Promoter) Promote(ctx context.Context, opts PromoteOptions) (*PromoteResult, error) {
	if _, err := p.registry.Require(opts.DatasetID); err != nil {
		return nil, err
	}

	lakeRoot := lake.NormalizeRoot(p.cfg.LakeLocalRoot)
	silverRoot := lake.NormalizeRoot(p.cfg.DeltaStoragePath)

	bronzeDir, err := lake.BronzeDir(lakeRoot, opts.DatasetID)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(bronzeDir); err != nil {
		return nil, fmt.Errorf("bronze directory missing (%s): run ingestor first", bronzeDir)
	}

	silverDir, err := lake.SilverTableDir(silverRoot, opts.DatasetID)
	if err != nil {
		return nil, err
	}
	if p.cfg.StorageMode == config.StorageModeLocal {
		if err := os.MkdirAll(filepath.Dir(silverDir), 0o755); err != nil {
			return nil, fmt.Errorf("create silver parent dir: %w", err)
		}
	}

	if _, err := os.Stat(p.script); err != nil {
		return nil, fmt.Errorf("promote script not found at %s", p.script)
	}

	cmd := exec.CommandContext(ctx, p.python, p.script,
		"--bronze-dir", bronzeDir,
		"--silver-dir", silverDir,
		"--dataset-id", opts.DatasetID,
		"--storage-mode", p.cfg.StorageMode,
		"--min-versions", fmt.Sprintf("%d", p.cfg.DeltaMinVersions),
	)
	cmd.Env = append(os.Environ(), p.storageEnv()...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("delta promote: %w\n%s", err, strings.TrimSpace(string(out)))
	}

	var payload struct {
		RowCount int `json:"row_count"`
	}
	if err := json.Unmarshal(lastJSONLine(out), &payload); err != nil {
		return nil, fmt.Errorf("parse promote output: %w (raw=%q)", err, string(out))
	}

	return &PromoteResult{
		DatasetID:   opts.DatasetID,
		RowCount:    payload.RowCount,
		BronzeDir:   bronzeDir,
		SilverDir:   silverDir,
		StorageMode: p.cfg.StorageMode,
	}, nil
}

func (p *Promoter) storageEnv() []string {
	switch p.cfg.StorageMode {
	case config.StorageModeMinIO:
		return []string{
			"MINIO_ENDPOINT=" + p.cfg.MinIOEndpoint,
			"MINIO_ACCESS_KEY=" + p.cfg.MinIOAccessKey,
			"MINIO_SECRET_KEY=" + p.cfg.MinIOSecretKey,
			"MINIO_BUCKET=" + p.cfg.MinIOBucket,
		}
	case config.StorageModeR2:
		return []string{
			"R2_ENDPOINT=" + p.cfg.R2Endpoint,
			"R2_ACCESS_KEY_ID=" + p.cfg.R2AccessKeyID,
			"R2_SECRET_ACCESS_KEY=" + p.cfg.R2SecretAccessKey,
			"R2_BUCKET=" + p.cfg.R2Bucket,
		}
	default:
		return nil
	}
}

func defaultPromoteScript() string {
	if v := strings.TrimSpace(os.Getenv("DELTA_PROMOTE_SCRIPT")); v != "" {
		return v
	}
	if root, err := findModuleRoot(); err == nil {
		return filepath.Join(root, "scripts", "delta", "promote.py")
	}
	return filepath.Join("scripts", "delta", "promote.py")
}

func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", dir)
		}
		dir = parent
	}
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func lastJSONLine(out []byte) []byte {
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "{") {
			return []byte(line)
		}
	}
	return out
}
