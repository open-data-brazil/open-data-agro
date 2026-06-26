package inmet

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFlattenSecaMonitorGolden(t *testing.T) {
	t.Parallel()

	raw := readINMETTestdata(t, "sequia_monitor.sample.json")
	headers, rows, err := FlattenSecaMonitor(raw)
	if err != nil {
		t.Fatalf("FlattenSecaMonitor: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows: got %d want 2", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["categoria_seca"]]; got != "S4" {
		t.Fatalf("categoria_seca: got %q want S4", got)
	}
	if got := rows[0][idx["ano"]]; got != "2026" {
		t.Fatalf("ano: got %q want 2026", got)
	}
}

func readINMETTestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
