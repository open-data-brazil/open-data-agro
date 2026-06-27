package ibge

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// PPMFetchOptions controls chunked SIDRA pulls for PPM datasets.
type PPMFetchOptions struct {
	FromYear int
	ToYear   int
	UFs      []string
}

// FetchPPMSnapshot downloads and merges SIDRA rows for a PPM catalog entry.
func (c *Client) FetchPPMSnapshot(ctx context.Context, entry catalog.RegistryEntry, opts PPMFetchOptions) ([]byte, string, error) {
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
		variables = "all"
	}

	ufChunks := chunkUFs(opts.UFs)
	var merged []map[string]any
	var requestURLs []string

	for year := from; year <= to; year++ {
		for _, ufBatch := range ufChunks {
			requestURL := buildPPMRequestURL(table, ufBatch, year, variables)
			requestURLs = append(requestURLs, requestURL)

			result, err := c.Download(ctx, requestURL)
			if err != nil {
				return nil, "", fmt.Errorf("sidra ppm fetch %d %s: %w", year, strings.Join(ufBatch, ","), err)
			}

			rows, err := parseSIDRARows(result.Body)
			if err != nil {
				return nil, "", fmt.Errorf("parse sidra ppm response %s: %w", requestURL, err)
			}
			merged = append(merged, rows...)
		}
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

func ppmUsesUFTerritory(table string) bool {
	switch table {
	case "3939", "94", "95", "3940":
		return true
	default:
		return false
	}
}

func buildPPMRequestURL(table string, ufBatch []string, year int, variables string) string {
	if ppmUsesUFTerritory(table) {
		return buildPPMUFURL(table, ufBatch, year, variables)
	}
	return buildPPMURL(table, ufBatch, year, variables)
}

func buildPPMURL(table string, ufBatch []string, year int, variables string) string {
	ufParam := strings.Join(ufBatch, ",")
	path := fmt.Sprintf("/t/%s/n6/in n3 %s/p/%d/v/%s", table, ufParam, year, variables)
	return sidraValuesBase + strings.ReplaceAll(path, " ", "%20")
}

func buildPPMUFURL(table string, ufBatch []string, year int, variables string) string {
	ufParam := strings.Join(ufBatch, ",")
	path := fmt.Sprintf("/t/%s/n3/%s/p/%d/v/%s", table, ufParam, year, variables)
	return sidraValuesBase + strings.ReplaceAll(path, " ", "%20")
}

// FlattenPPMUF converts UF-level PPM SIDRA JSON into bronze columns.
func FlattenPPMUF(datasetID string, raw []byte) ([]string, [][]string, error) {
	var rows []map[string]any
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse ppm uf json: %w", err)
	}

	table := sidraTableForPPMDataset(datasetID)
	headers := []string{
		"sidra_tabela",
		"codigo_uf",
		"uf",
		"ano",
		"variavel_codigo",
		"variavel",
		"categoria_codigo",
		"categoria",
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
		return nil, nil, fmt.Errorf("no ppm uf rows to flatten for %s", datasetID)
	}
	return headers, out, nil
}

func sidraTableForPPMDataset(datasetID string) string {
	switch datasetID {
	case "ibge.ppm-producao-municipal":
		return "74"
	case "ibge.ppm-efetivo-rebanhos":
		return "3939"
	case "ibge.ppm-vacas-ordenhadas":
		return "94"
	case "ibge.ppm-ovinos-tosquiados":
		return "95"
	case "ibge.ppm-aquicultura":
		return "3940"
	default:
		return ""
	}
}

func isPPMUFDataset(datasetID string) bool {
	switch datasetID {
	case "ibge.ppm-efetivo-rebanhos", "ibge.ppm-vacas-ordenhadas", "ibge.ppm-ovinos-tosquiados", "ibge.ppm-aquicultura":
		return true
	default:
		return false
	}
}

// ResolvePPMURL validates the catalog base URL for a PPM dataset.
func ResolvePPMURL(entry catalog.RegistryEntry) (string, error) {
	return ResolvePAMURL(entry)
}

func isPPMDataset(datasetID string) bool {
	return strings.HasPrefix(datasetID, "ibge.ppm-")
}
