package anp

import "testing"

func TestFilterEthanolPrecos(t *testing.T) {
	t.Parallel()

	headers := []string{"DATA INICIAL", "PRODUTO", "PREÇO MÉDIO REVENDA"}
	rows := [][]string{
		{"2024-01-01", "ETANOL HIDRATADO", "3.50"},
		{"2024-01-01", "GASOLINA COMUM", "5.80"},
		{"2024-01-01", "ETANOL ANIDRO", "4.10"},
	}

	_, filtered := FilterEthanolPrecos(headers, rows)
	if len(filtered) != 2 {
		t.Fatalf("filtered rows: got %d want 2", len(filtered))
	}
	if filtered[0][1] != "ETANOL HIDRATADO" {
		t.Fatalf("first product: got %q", filtered[0][1])
	}
}
