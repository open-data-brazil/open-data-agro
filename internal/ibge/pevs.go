package ibge

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const pevsYearChunkSize = 10

// PEVSFetchOptions controls chunked SIDRA pulls for PEVS datasets.
type PEVSFetchOptions struct {
	FromYear int
	ToYear   int
}

// FetchPEVSSnapshot downloads and merges SIDRA rows for a PEVS catalog entry.
func (c *Client) FetchPEVSSnapshot(ctx context.Context, entry catalog.RegistryEntry, opts PEVSFetchOptions) ([]byte, string, error) {
	table := strings.TrimSpace(entry.SidraTable)
	if table == "" {
		return nil, "", fmt.Errorf("dataset %s missing sidra_table", entry.DatasetID)
	}

	from, to, err := resolveYearRange(entry, opts.FromYear, opts.ToYear)
	if err != nil {
		return nil, "", err
	}

	variables := formatVariables(entry.SidraVariables)
	if variables == "" {
		return nil, "", fmt.Errorf("dataset %s missing sidra_variables", entry.DatasetID)
	}

	var merged []map[string]any
	var requestURLs []string

	for start := from; start <= to; start += pevsYearChunkSize {
		end := start + pevsYearChunkSize - 1
		if end > to {
			end = to
		}
		years := formatYearRange(start, end)
		requestURL := buildPEVSURL(table, years, variables)
		requestURLs = append(requestURLs, requestURL)

		result, err := c.Download(ctx, requestURL)
		if err != nil {
			return nil, "", fmt.Errorf("sidra pevs fetch %s: %w", years, err)
		}

		rows, err := parseSIDRARows(result.Body)
		if err != nil {
			return nil, "", fmt.Errorf("parse sidra pevs response %s: %w", requestURL, err)
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

func buildPEVSURL(table, years, variables string) string {
	path := fmt.Sprintf("/t/%s/n3/all/p/%s/v/%s", table, years, variables)
	return sidraValuesBase + strings.ReplaceAll(path, " ", "%20")
}

func formatYearRange(from, to int) string {
	if from == to {
		return strconv.Itoa(from)
	}
	years := make([]string, 0, to-from+1)
	for year := from; year <= to; year++ {
		years = append(years, strconv.Itoa(year))
	}
	return strings.Join(years, ",")
}

// FlattenPEVS converts merged SIDRA PEVS JSON rows into canonical bronze columns.
func FlattenPEVS(datasetID string, raw []byte) ([]string, [][]string, error) {
	var rows []map[string]any
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse pevs json: %w", err)
	}

	table := sidraTableForPEVSDataset(datasetID)
	headers := []string{
		"sidra_tabela",
		"codigo_uf",
		"uf",
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
		return nil, nil, fmt.Errorf("no pevs rows to flatten for %s", datasetID)
	}
	return headers, out, nil
}

func sidraTableForPEVSDataset(datasetID string) string {
	switch datasetID {
	case "ibge.pevs-producao-vegetal":
		return "289"
	default:
		return ""
	}
}

// ResolvePEVSURL validates the catalog base URL for a PEVS dataset.
func ResolvePEVSURL(entry catalog.RegistryEntry) (string, error) {
	return ResolvePAMURL(entry)
}

func isPEVSDataset(datasetID string) bool {
	return strings.HasPrefix(datasetID, "ibge.pevs-")
}
