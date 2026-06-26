package catalog

import (
	"testing"
	"time"
)

func TestParseDatasetID(t *testing.T) {
	t.Parallel()

	valid := []string{
		"conab.estimativa-graos",
		"conab.precos-agropecuarios-mensal-municipio",
		"oecd-fao.ag-outlook",
	}
	for _, id := range valid {
		if _, err := ParseDatasetID(id); err != nil {
			t.Fatalf("expected valid ID %q: %v", id, err)
		}
	}

	invalid := []string{"", "CONAB.foo", "conab", "conab.", ".foo.bar"}
	for _, id := range invalid {
		if _, err := ParseDatasetID(id); err == nil {
			t.Fatalf("expected invalid ID %q", id)
		}
	}
}

func TestLoadRegistryFromConabYAML(t *testing.T) {
	t.Parallel()

	root, err := repoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}

	reg, err := LoadRegistryFrom(root + "/configs/catalog")
	if err != nil {
		t.Fatalf("load registry: %v", err)
	}

	entry, ok := reg.Get("conab.estimativa-graos")
	if !ok {
		t.Fatal("expected conab.estimativa-graos in registry")
	}
	if entry.Format != FormatTXT {
		t.Fatalf("format: got %q want %q", entry.Format, FormatTXT)
	}
	if entry.Delimiter != ";" {
		t.Fatalf("delimiter: got %q want %q", entry.Delimiter, ";")
	}
	if entry.DiscoveredAt.IsZero() {
		t.Fatal("expected discovered_at to be set")
	}
	if entry.DiscoveredAt.UTC() != time.Date(2026, 6, 25, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("discovered_at: got %v", entry.DiscoveredAt.UTC())
	}

	ids := reg.ListIDs()
	if len(ids) < 20 {
		t.Fatalf("expected at least 20 CONAB datasets, got %d", len(ids))
	}
}

func TestRegistryRequireUnknown(t *testing.T) {
	t.Parallel()

	reg := NewRegistry(nil)
	if _, err := reg.Require("conab.missing"); err == nil {
		t.Fatal("expected error for unknown dataset")
	}
}
