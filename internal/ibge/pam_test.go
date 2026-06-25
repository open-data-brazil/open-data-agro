package ibge

import (
	"testing"
)

func TestFlattenPAMSojaSample(t *testing.T) {
	t.Parallel()

	raw := readTestdata(t, "pam_area_quantidade.sample.json")
	headers, rows, err := FlattenPAM("ibge.pam-area-quantidade", raw)
	if err != nil {
		t.Fatalf("FlattenPAM: %v", err)
	}
	if len(rows) != 6 {
		t.Fatalf("rows: got %d want 6", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}

	row := rows[0]
	if got := row[idx["sidra_tabela"]]; got != "1612" {
		t.Fatalf("sidra_tabela: got %q want 1612", got)
	}
	if got := row[idx["codigo_ibge"]]; got != "4300034" {
		t.Fatalf("codigo_ibge: got %q want 4300034", got)
	}
	if got := row[idx["ano"]]; got != "2015" {
		t.Fatalf("ano: got %q want 2015", got)
	}
	if got := row[idx["produto_codigo"]]; got != "2713" {
		t.Fatalf("produto_codigo: got %q want 2713", got)
	}
	if got := row[idx["variavel_codigo"]]; got != "109" {
		t.Fatalf("variavel_codigo: got %q want 109", got)
	}
}
