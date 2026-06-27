package inpe

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFlattenDETERGolden(t *testing.T) {
	t.Parallel()

	raw := readINPETestdata(t, "deter_alertas_desmatamento.sample.json")
	headers, rows, err := FlattenDETER(raw)
	if err != nil {
		t.Fatalf("FlattenDETER: %v", err)
	}
	if len(rows) != 5 {
		t.Fatalf("rows: got %d want 5", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["uf"]]; got != "PA" {
		t.Fatalf("uf: got %q want PA", got)
	}
	if got := rows[0][idx["class_name"]]; got == "" {
		t.Fatalf("class_name should not be empty")
	}
}

func readINPETestdata(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
