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
		return nil, "", fmt.Errorf("dataset %s missing sidra_variables", entry.DatasetID)
	}

	ufChunks := chunkUFs(opts.UFs)
	var merged []map[string]any
	var requestURLs []string

	for year := from; year <= to; year++ {
		for _, ufBatch := range ufChunks {
			requestURL := buildPPMURL(table, ufBatch, year, variables)
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

func buildPPMURL(table string, ufBatch []string, year int, variables string) string {
	ufParam := strings.Join(ufBatch, ",")
	path := fmt.Sprintf("/t/%s/n6/in n3 %s/p/%d/v/%s", table, ufParam, year, variables)
	return sidraValuesBase + strings.ReplaceAll(path, " ", "%20")
}

// ResolvePPMURL validates the catalog base URL for a PPM dataset.
func ResolvePPMURL(entry catalog.RegistryEntry) (string, error) {
	return ResolvePAMURL(entry)
}

func isPPMDataset(datasetID string) bool {
	return strings.HasPrefix(datasetID, "ibge.ppm-")
}
