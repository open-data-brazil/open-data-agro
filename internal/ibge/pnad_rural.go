package ibge

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const pnadQuarterChunkSize = 4

// PNADFetchOptions controls chunked SIDRA pulls for PNAD datasets.
type PNADFetchOptions struct {
	UFs []string
}

// FetchPNADRuralSnapshot downloads PNAD Contínua labor rows for catalog entry.
func (c *Client) FetchPNADRuralSnapshot(ctx context.Context, entry catalog.RegistryEntry, opts PNADFetchOptions) ([]byte, string, error) {
	table := strings.TrimSpace(entry.SidraTable)
	if table == "" {
		return nil, "", fmt.Errorf("dataset %s missing sidra_table", entry.DatasetID)
	}

	variables := formatVariables(entry.SidraVariables)
	if variables == "" {
		variables = "all"
	}

	ufChunks := chunkUFs(opts.UFs)
	if len(ufChunks) == 1 && len(opts.UFs) == 0 {
		ufChunks = defaultUFChunks
	}

	var merged []map[string]any
	var requestURLs []string

	for _, ufBatch := range ufChunks {
		requestURL := buildPNADURL(table, ufBatch, "last 12", variables)
		requestURLs = append(requestURLs, requestURL)

		result, err := c.Download(ctx, requestURL)
		if err != nil {
			return nil, "", fmt.Errorf("sidra pnad fetch %s: %w", strings.Join(ufBatch, ","), err)
		}

		rows, err := parseSIDRARows(result.Body)
		if err != nil {
			return nil, "", fmt.Errorf("parse sidra pnad response %s: %w", requestURL, err)
		}
		merged = append(merged, rows...)
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("sidra returned no data rows for %s", entry.DatasetID)
	}

	payload, err := json.Marshal(merged)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s/t/%s (chunks: %d)", sidraValuesBase, table, len(requestURLs))
	return payload, sourceURL, nil
}

func buildPNADURL(table string, ufBatch []string, period, variables string) string {
	ufParam := strings.Join(ufBatch, ",")
	path := fmt.Sprintf("/t/%s/n3/%s/p/%s/v/%s", table, ufParam, period, variables)
	return sidraValuesBase + strings.ReplaceAll(path, " ", "%20")
}

// FlattenPNADRural converts merged SIDRA PNAD JSON into bronze columns.
func FlattenPNADRural(datasetID string, raw []byte) ([]string, [][]string, error) {
	var rows []map[string]any
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse pnad json: %w", err)
	}

	table := sidraTableForPNADDataset(datasetID)
	headers := []string{
		"sidra_tabela",
		"codigo_uf",
		"uf",
		"trimestre",
		"variavel_codigo",
		"variavel",
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
			stringField(row["D2N"]),
			formatID(row["D3C"]),
			stringField(row["D3N"]),
			stringField(row["V"]),
			formatID(row["MC"]),
			stringField(row["MN"]),
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no pnad rows to flatten for %s", datasetID)
	}
	return headers, out, nil
}

func sidraTableForPNADDataset(datasetID string) string {
	switch datasetID {
	case "ibge.pnad-continua-rural":
		return "6385"
	default:
		return ""
	}
}

// ResolvePNADRuralURL validates the catalog base URL for a PNAD dataset.
func ResolvePNADRuralURL(entry catalog.RegistryEntry) (string, error) {
	return ResolvePAMURL(entry)
}

func isPNADRuralDataset(datasetID string) bool {
	return strings.HasPrefix(datasetID, "ibge.pnad-continua-")
}
