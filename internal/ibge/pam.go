package ibge

import (
	"encoding/json"
	"fmt"
)

// FlattenPAM converts merged SIDRA PAM JSON rows into canonical bronze columns.
func FlattenPAM(datasetID string, raw []byte) ([]string, [][]string, error) {
	var rows []map[string]any
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse pam json: %w", err)
	}

	table := sidraTableForDataset(datasetID)
	headers := []string{
		"sidra_tabela",
		"codigo_ibge",
		"municipio",
		"ano",
		"variavel_codigo",
		"variavel",
		"produto_codigo",
		"produto",
		"valor",
		"unidade_codigo",
		"unidade",
	}

	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			table,
			formatID(row["D1C"]),
			stringField(row["D1N"]),
			formatID(row["D2C"]),
			formatID(row["D3C"]),
			stringField(row["D3N"]),
			formatID(row["D4C"]),
			stringField(row["D4N"]),
			stringField(row["V"]),
			formatID(row["MC"]),
			stringField(row["MN"]),
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no pam rows to flatten for %s", datasetID)
	}
	return headers, out, nil
}

func sidraTableForDataset(datasetID string) string {
	switch datasetID {
	case "ibge.pam-area-quantidade":
		return "1612"
	case "ibge.pam-rendimento-valor":
		return "1613"
	case "ibge.pam-estabelecimentos":
		return "5457"
	default:
		return ""
	}
}

// FlattenIBGEJSON routes IBGE JSON payloads to the correct flattener.
func FlattenIBGEJSON(datasetID string, raw []byte) ([]string, [][]string, error) {
	if isPAMDataset(datasetID) {
		return FlattenPAM(datasetID, raw)
	}
	if isLSPADataset(datasetID) {
		return FlattenLSPA(datasetID, raw)
	}
	return FlattenLocalidades(datasetID, raw)
}
