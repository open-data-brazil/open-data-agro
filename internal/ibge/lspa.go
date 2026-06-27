package ibge

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// LSPACropSlug maps LSPA SIDRA crop codes to canonical slugs.
var LSPACropSlug = map[string]string{
	"39443": "soja",
	"39441": "milho",
	"39445": "trigo",
}

// LSPAFetchOptions controls chunked SIDRA pulls for LSPA datasets.
type LSPAFetchOptions struct {
	Crop     string
	FromYear int
	ToYear   int
	UFs      []string
}

// FetchLSPASnapshot downloads and merges SIDRA rows for an LSPA catalog entry.
func (c *Client) FetchLSPASnapshot(ctx context.Context, entry catalog.RegistryEntry, opts LSPAFetchOptions) ([]byte, string, error) {
	table := strings.TrimSpace(entry.SidraTable)
	if table == "" {
		return nil, "", fmt.Errorf("dataset %s missing sidra_table", entry.DatasetID)
	}

	classification := strings.TrimSpace(entry.SidraClassification)
	if classification == "" {
		return nil, "", fmt.Errorf("dataset %s missing sidra_classification", entry.DatasetID)
	}

	crops, err := resolveCropCodes(entry, opts.Crop)
	if err != nil {
		return nil, "", err
	}

	from, to, err := resolveYearRange(entry, opts.FromYear, opts.ToYear)
	if err != nil {
		return nil, "", err
	}

	variables := formatVariables(entry.SidraVariables)
	if variables == "" {
		return nil, "", fmt.Errorf("dataset %s missing sidra_variables", entry.DatasetID)
	}

	ufChunks := chunkUFs(opts.UFs)
	var merged []map[string]any
	var requestURLs []string

	for year := from; year <= to; year++ {
		periods := formatYearMonths(year)
		for cropName, cropCode := range crops {
			for _, ufBatch := range ufChunks {
				requestURL := buildLSPAURL(table, ufBatch, periods, variables, classification, cropCode)
				requestURLs = append(requestURLs, requestURL)

				result, err := c.Download(ctx, requestURL)
				if err != nil {
					return nil, "", fmt.Errorf("sidra lspa fetch %s %d %s: %w", cropName, year, strings.Join(ufBatch, ","), err)
				}

				rows, err := parseSIDRARows(result.Body)
				if err != nil {
					return nil, "", fmt.Errorf("parse sidra lspa response %s: %w", requestURL, err)
				}
				merged = append(merged, rows...)
			}
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

func buildLSPAURL(table string, ufBatch []string, periods, variables, classification, cropCode string) string {
	ufParam := strings.Join(ufBatch, ",")
	path := fmt.Sprintf(
		"/t/%s/n3/in %s/p/%s/v/%s/c%s/%s",
		table,
		ufParam,
		periods,
		variables,
		classification,
		cropCode,
	)
	return sidraValuesBase + strings.ReplaceAll(path, " ", "%20")
}

func formatYearMonths(year int) string {
	months := make([]string, 12)
	for month := 1; month <= 12; month++ {
		months[month-1] = fmt.Sprintf("%d%02d", year, month)
	}
	return strings.Join(months, ",")
}

// FlattenLSPA converts merged SIDRA LSPA JSON rows into canonical bronze columns.
func FlattenLSPA(datasetID string, raw []byte) ([]string, [][]string, error) {
	var rows []map[string]any
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse lspa json: %w", err)
	}

	table := sidraTableForLSPADataset(datasetID)
	headers := []string{
		"sidra_tabela",
		"codigo_uf",
		"uf",
		"mes",
		"variavel_codigo",
		"variavel",
		"produto_codigo",
		"produto",
		"produto_slug",
		"valor",
		"unidade_codigo",
		"unidade",
	}

	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		produtoCode := formatID(row["D4C"])
		slug := LSPACropSlug[produtoCode]
		if slug == "" {
			slug = slugFromCropName(stringField(row["D4N"]))
		}
		out = append(out, []string{
			table,
			formatID(row["D1C"]),
			stringField(row["D1N"]),
			formatID(row["D2C"]),
			formatID(row["D3C"]),
			stringField(row["D3N"]),
			produtoCode,
			stringField(row["D4N"]),
			slug,
			stringField(row["V"]),
			formatID(row["MC"]),
			stringField(row["MN"]),
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no lspa rows to flatten for %s", datasetID)
	}
	return headers, out, nil
}

func sidraTableForLSPADataset(datasetID string) string {
	switch datasetID {
	case "ibge.lspa-area-producao":
		return "6588"
	case "ibge.lspa-rendimento-medio":
		return "6588"
	default:
		return ""
	}
}

func slugFromCropName(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	switch {
	case strings.Contains(lower, "soja"):
		return "soja"
	case strings.Contains(lower, "milho"):
		return "milho"
	case strings.Contains(lower, "trigo"):
		return "trigo"
	default:
		return ""
	}
}

// ResolveLSPAURL validates the catalog base URL for an LSPA dataset.
func ResolveLSPAURL(entry catalog.RegistryEntry) (string, error) {
	return ResolvePAMURL(entry)
}

func isLSPADataset(datasetID string) bool {
	return strings.HasPrefix(datasetID, "ibge.lspa-")
}
