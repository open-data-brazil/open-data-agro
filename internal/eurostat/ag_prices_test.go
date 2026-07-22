package eurostat

import (
	"os"
	"strings"
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

func TestBuildDatasetURLApriPiOuta(t *testing.T) {
	t.Parallel()

	got := buildDatasetURL("apri_pi_outa", "EU27_2020", []string{"010000", "AM011000"}, 2020)
	for _, want := range []string{
		"/apri_pi_outa?",
		"am_item=AM010000",
		"am_item=AM011000",
		"unit=I15",
		"p_adj=NI",
		"geo=EU27_2020",
		"sinceTimePeriod=2020",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in %q", want, got)
		}
	}
	if strings.Contains(got, "product=") {
		t.Fatalf("legacy product param must not appear: %q", got)
	}
}
