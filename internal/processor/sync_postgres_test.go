package processor

import "testing"

func TestMartTableName(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"mart_conab__estimativa_graos":        "conab_estimativa_graos",
		"mart_ibge__localidades_municipios":   "ibge_localidades_municipios",
		"mart_noaa__global_temp_anomaly":      "noaa_global_temp_anomaly",
		"mart_worldbank__pink_sheet_monthly":  "worldbank_pink_sheet_monthly",
	}

	for dir, want := range cases {
		got, err := MartTableName(dir)
		if err != nil {
			t.Fatalf("MartTableName(%q): %v", dir, err)
		}
		if got != want {
			t.Fatalf("MartTableName(%q) = %q, want %q", dir, got, want)
		}
	}
}

func TestDateRangeFromRows(t *testing.T) {
	t.Parallel()

	columns := []string{"refmonth", "value"}
	rows := [][]string{
		{"2024-02", "1"},
		{"2024-01", "2"},
	}
	minDate, maxDate := dateRangeFromRows(columns, rows)
	if minDate != "2024-01" || maxDate != "2024-02" {
		t.Fatalf("range: got %q-%q", minDate, maxDate)
	}
}

func TestParseMartFilter(t *testing.T) {
	t.Parallel()

	got := ParseMartFilter(" conab_estimativa_graos, ibge_localidades_municipios ")
	if len(got) != 2 || got[0] != "conab_estimativa_graos" {
		t.Fatalf("filter: %#v", got)
	}
}
