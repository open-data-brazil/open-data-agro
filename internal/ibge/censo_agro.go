package ibge

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// FetchCensoAgroSnapshot downloads SIDRA Censo Agro 2017 establishment rows by UF.
func (c *Client) FetchCensoAgroSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	table := strings.TrimSpace(entry.SidraTable)
	if table == "" {
		return nil, "", fmt.Errorf("dataset %s missing sidra_table", entry.DatasetID)
	}

	year := entry.PeriodEnd
	if year == 0 {
		year = 2017
	}

	variables := formatVariables(entry.SidraVariables)
	path := fmt.Sprintf("/t/%s/n3/all/p/%d", table, year)
	if variables != "" {
		path += "/v/" + variables
	} else {
		path += "/v/all"
	}
	requestURL := sidraValuesBase + path

	result, err := c.Download(ctx, requestURL)
	if err != nil {
		return nil, "", fmt.Errorf("sidra censo agro fetch: %w", err)
	}

	rows, err := parseSIDRARows(result.Body)
	if err != nil {
		return nil, "", fmt.Errorf("parse sidra censo agro response: %w", err)
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("sidra returned no data rows for %s", entry.DatasetID)
	}

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s/t/%s (year: %d)", sidraValuesBase, table, year)
	return payload, sourceURL, nil
}

// FlattenCensoAgro converts merged SIDRA Censo Agro JSON into bronze columns.
func FlattenCensoAgro(datasetID string, raw []byte) ([]string, [][]string, error) {
	var rows []map[string]any
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse censo agro json: %w", err)
	}

	table := sidraTableForCensoAgroDataset(datasetID)
	headers := []string{
		"sidra_tabela",
		"codigo_uf",
		"uf",
		"ano",
		"variavel_codigo",
		"variavel",
		"condicao_produtor_codigo",
		"condicao_produtor",
		"tipologia_codigo",
		"tipologia",
		"atividade_codigo",
		"atividade",
		"sexo_produtor_codigo",
		"sexo_produtor",
		"idade_produtor_codigo",
		"idade_produtor",
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
			formatID(row["D5C"]),
			stringField(row["D5N"]),
			formatID(row["D6C"]),
			stringField(row["D6N"]),
			formatID(row["D7C"]),
			stringField(row["D7N"]),
			formatID(row["D8C"]),
			stringField(row["D8N"]),
			stringField(row["V"]),
			formatID(row["MC"]),
			stringField(row["MN"]),
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no censo agro rows to flatten for %s", datasetID)
	}
	return headers, out, nil
}

func sidraTableForCensoAgroDataset(datasetID string) string {
	switch datasetID {
	case "ibge.censo-agro-estabelecimentos":
		return "6878"
	case "ibge.censo-agro-area-uso-solo":
		return "6879"
	case "ibge.censo-agro-maquinario":
		return "6880"
	default:
		return ""
	}
}

// ResolveCensoAgroURL validates the catalog base URL for a Censo Agro dataset.
func ResolveCensoAgroURL(entry catalog.RegistryEntry) (string, error) {
	return ResolvePAMURL(entry)
}

func isCensoAgroDataset(datasetID string) bool {
	return strings.HasPrefix(datasetID, "ibge.censo-agro-")
}
