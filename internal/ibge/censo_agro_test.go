package ibge

import (
	"testing"
)

func TestFlattenCensoAgroGolden(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "censo_agro_estabelecimentos.sample.json")
	headers, rows, err := FlattenCensoAgro("ibge.censo-agro-estabelecimentos", raw)
	if err != nil {
		t.Fatalf("FlattenCensoAgro: %v", err)
	}
	if len(rows) != 55 {
		t.Fatalf("rows: got %d want 55", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["sidra_tabela"]]; got != "6878" {
		t.Fatalf("sidra_tabela: got %q want 6878", got)
	}
	if got := rows[0][idx["codigo_uf"]]; got == "" {
		t.Fatalf("codigo_uf should not be empty")
	}
}

func TestFlattenPNADRuralGolden(t *testing.T) {
	t.Parallel()

	raw := readIBGETestdata(t, "pnad_continua_rural.sample.json")
	headers, rows, err := FlattenPNADRural("ibge.pnad-continua-rural", raw)
	if err != nil {
		t.Fatalf("FlattenPNADRural: %v", err)
	}
	if len(rows) != 33 {
		t.Fatalf("rows: got %d want 33", len(rows))
	}

	idx := map[string]int{}
	for i, h := range headers {
		idx[h] = i
	}
	if got := rows[0][idx["sidra_tabela"]]; got != "6385" {
		t.Fatalf("sidra_tabela: got %q want 6385", got)
	}
}
