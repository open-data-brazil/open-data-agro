package processor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/config"
)

func writeCONABBronze(t *testing.T, root string, nullProduto bool) {
	t.Helper()

	bronzeDir := filepath.Join(root, "bronze", "conab", "estimativa-graos", "ingest_date=2026-06-25")
	if err := os.MkdirAll(bronzeDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	script := defaultQualityScript()
	if _, err := os.Stat(script); err != nil {
		t.Skipf("quality script missing at %s", script)
	}

	python := envOr("PYTHON", "python3")
	writeScript := `
import sys
from pathlib import Path
import pyarrow as pa
import pyarrow.parquet as pq
out = Path(sys.argv[1])
null_produto = sys.argv[2] == "true"
produto = [None] if null_produto else ["Soja"]
pq.write_table(pa.table({
    "produto": produto,
    "uf": ["PR"],
    "safra": ["UNICA"],
    "ano_agricola": ["2025/26"],
}), out / "part-test.parquet")
`
	cmd := exec.Command(python, "-c", writeScript, bronzeDir, fmt.Sprintf("%t", nullProduto))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("write bronze parquet: %v\n%s", err, out)
	}
}

func TestQualityRequiresBronze(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	cfg := config.LakeConfig{
		StorageMode:   config.StorageModeLocal,
		LakeLocalRoot: root,
	}
	reg := testRegistry(catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	})

	gate := NewQualityGate(cfg, reg)
	_, err := gate.RunBronzeCheckpoint(context.Background(), QualityOptions{DatasetID: "conab.estimativa-graos"})
	if err == nil {
		t.Fatal("expected error when bronze is missing")
	}
}

func TestQualityLocalIntegration(t *testing.T) {
	if os.Getenv("GE_INTEGRATION") != "1" {
		t.Skip("set GE_INTEGRATION=1 with great-expectations installed")
	}

	root := t.TempDir()
	cfg := config.LakeConfig{
		StorageMode:   config.StorageModeLocal,
		LakeLocalRoot: root,
	}
	reg := testRegistry(catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	})

	writeCONABBronze(t, root, false)

	gate := NewQualityGate(cfg, reg)
	result, err := gate.RunBronzeCheckpoint(context.Background(), QualityOptions{DatasetID: "conab.estimativa-graos"})
	if err != nil {
		t.Fatalf("RunBronzeCheckpoint: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected success, got %+v", result)
	}
}

func TestQualityFailsOnBadBronze(t *testing.T) {
	if os.Getenv("GE_INTEGRATION") != "1" {
		t.Skip("set GE_INTEGRATION=1 with great-expectations installed")
	}

	root := t.TempDir()
	cfg := config.LakeConfig{
		StorageMode:   config.StorageModeLocal,
		LakeLocalRoot: root,
	}
	reg := testRegistry(catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("conab.estimativa-graos"),
		Format:    catalog.FormatTXT,
		Delimiter: ";",
	})

	writeCONABBronze(t, root, true)

	gate := NewQualityGate(cfg, reg)
	_, err := gate.RunBronzeCheckpoint(context.Background(), QualityOptions{DatasetID: "conab.estimativa-graos"})
	if err == nil {
		t.Fatal("expected quality failure for null Produto")
	}
}
