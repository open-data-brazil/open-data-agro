package ana

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHidrologiaXML(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("testdata", "hidrologia_series.sample.xml"))
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}

	rows, err := parseHidrologiaXML(raw, "3")
	if err != nil {
		t.Fatalf("parseHidrologiaXML: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows: got %d want 2", len(rows))
	}
	if rows[0].StationCode != "15400000" {
		t.Fatalf("station: got %q", rows[0].StationCode)
	}
}
