package processor

import "testing"

func TestSubstituteVars(t *testing.T) {
	t.Parallel()

	sql := "SELECT * FROM read_parquet('${bronze_uri}');"
	got := SubstituteVars(sql, map[string]string{"bronze_uri": "lake/bronze/x/**/*.parquet"})
	want := "SELECT * FROM read_parquet('lake/bronze/x/**/*.parquet');"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestParseCountCSV(t *testing.T) {
	t.Parallel()

	count, err := parseCountCSV("row_count\n42\n")
	if err != nil {
		t.Fatalf("parseCountCSV: %v", err)
	}
	if count != 42 {
		t.Fatalf("got %d want 42", count)
	}
}
