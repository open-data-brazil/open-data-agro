package ibge

import (
	"testing"
)

func TestFlattenPPMGolden(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "ppm_producao_municipal.sample.json")
	headers, rows, err := FlattenPAM("ibge.ppm-producao-municipal", raw)
	if err != nil {
		t.Fatalf("FlattenPAM: %v", err)
	}
	if len(rows) != 7 {
		t.Fatalf("rows: got %d want 7", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["sidra_tabela"]]; got != "74" {
		t.Fatalf("sidra_tabela: got %q want 74", got)
	}
	if got := rows[0][idx["codigo_ibge"]]; got == "" {
		t.Fatalf("codigo_ibge should not be empty")
	}
}

func TestBuildPPMURL(t *testing.T) {
	t.Parallel()

	got := buildPPMURL("74", []string{"11", "12"}, 2023, "106,215")
	want := "https://apisidra.ibge.gov.br/values/t/74/n6/in%20n3%2011,12/p/2023/v/106,215"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
