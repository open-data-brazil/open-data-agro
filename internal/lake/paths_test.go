package lake

import (
	"path/filepath"
	"testing"
)

func TestSilverTableName(t *testing.T) {
	t.Parallel()

	name, err := SilverTableName("conab.estimativa-graos")
	if err != nil {
		t.Fatalf("SilverTableName: %v", err)
	}
	if name != "estimativa_graos" {
		t.Fatalf("got %q want estimativa_graos", name)
	}
}

func TestBronzeGlob(t *testing.T) {
	t.Parallel()

	glob, err := BronzeGlob("./lake", "conab.estimativa-graos")
	if err != nil {
		t.Fatalf("BronzeGlob: %v", err)
	}
	want := filepath.Join("lake", "bronze", "conab", "estimativa-graos", "**", "*.parquet")
	if glob != want {
		t.Fatalf("got %q want %q", glob, want)
	}
}

func TestSilverTableDir(t *testing.T) {
	t.Parallel()

	dir, err := SilverTableDir("./lake/silver", "conab.estimativa-graos")
	if err != nil {
		t.Fatalf("SilverTableDir: %v", err)
	}
	want := filepath.Join("lake", "silver", "conab", "estimativa_graos")
	if dir != want {
		t.Fatalf("got %q want %q", dir, want)
	}
}
