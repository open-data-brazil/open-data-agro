package mapa

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestParseSIPEAGROCSVGolden(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("testdata", "sipeagro_fertilizantes.sample.csv"))
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}

	rows, err := parseSIPEAGROCSV(raw, "Fertilizantes")
	if err != nil {
		t.Fatalf("parseSIPEAGROCSV: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows: got %d want 3", len(rows))
	}
	if rows[0][0] != "Fertilizantes" {
		t.Fatalf("linha_produto: got %q", rows[0][0])
	}
	if rows[0][1] != "SC" {
		t.Fatalf("uf: got %q want SC", rows[0][1])
	}
	if rows[0][3] == "" {
		t.Fatal("expected numero_registro_estabelecimento")
	}
}

func TestMergeSIPEAGROProdutosSample(t *testing.T) {
	t.Parallel()

	samples := map[string][]byte{
		"Fertilizantes":      mustRead(t, "sipeagro_fertilizantes.sample.csv"),
		"Qualidade Vegetal":  mustRead(t, "sipeagro_qualidade_vegetal.sample.csv"),
	}
	merged, err := MergeSIPEAGROSampleCSV("mapa.sipeagro-produtos", samples)
	if err != nil {
		t.Fatalf("MergeSIPEAGROSampleCSV: %v", err)
	}
	rows, err := parseSIPEAGROCSV(merged, "ignored")
	if err != nil {
		t.Fatalf("parse merged: %v", err)
	}
	if len(rows) != 6 {
		t.Fatalf("merged rows: got %d want 6", len(rows))
	}
}

func TestResolveSIPEAGROURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("mapa.sipeagro-estabelecimentos"),
	}
	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" || !containsAll(url, "package_show", "sipeagro") {
		t.Fatalf("unexpected url: %q", url)
	}
}

func mustRead(t *testing.T, name string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return raw
}

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !contains(value, part) {
			return false
		}
	}
	return true
}

func contains(value, part string) bool {
	return len(part) == 0 || (len(value) >= len(part) && indexOf(value, part) >= 0)
}

func indexOf(value, part string) int {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}
