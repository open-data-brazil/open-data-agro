package eurostat

import (
	"os"
	"testing"
)

func TestParseAgPricesJSON(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("testdata/apri_pi15_outa.sample.json")
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}

	rows, err := parseAgPricesJSON(raw, "apri_pi15_outa", "EU27_2020")
	if err != nil {
		t.Fatalf("parseAgPricesJSON: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("expected rows")
	}
	if rows[0].ProductCode == "" || rows[0].Year == "" {
		t.Fatalf("unexpected row: %+v", rows[0])
	}
}
