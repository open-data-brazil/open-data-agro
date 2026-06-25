package ibge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const sidraValuesBase = "https://apisidra.ibge.gov.br/values"

// Brazilian UF numeric codes (IBGE).
var defaultUFChunks = [][]string{
	{"11", "12", "13", "14", "15", "16", "17"},
	{"21", "22", "24", "25", "26", "27", "28", "29"},
	{"31", "32", "33", "35", "41", "42", "43"},
	{"50", "51", "52", "53"},
}

// PAMFetchOptions controls chunked SIDRA pulls for PAM datasets.
type PAMFetchOptions struct {
	Crop     string
	FromYear int
	ToYear   int
	UFs      []string
}

// FetchPAMSnapshot downloads and merges SIDRA rows for a PAM catalog entry.
func (c *Client) FetchPAMSnapshot(ctx context.Context, entry catalog.RegistryEntry, opts PAMFetchOptions) ([]byte, string, error) {
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
		for cropName, cropCode := range crops {
			for _, ufBatch := range ufChunks {
				requestURL := buildSIDRAURL(table, ufBatch, year, variables, classification, cropCode)
				requestURLs = append(requestURLs, requestURL)

				result, err := c.Download(ctx, requestURL)
				if err != nil {
					return nil, "", fmt.Errorf("sidra fetch %s %d %s: %w", cropName, year, strings.Join(ufBatch, ","), err)
				}

				rows, err := parseSIDRARows(result.Body)
				if err != nil {
					return nil, "", fmt.Errorf("parse sidra response %s: %w", requestURL, err)
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

func buildSIDRAURL(table string, ufBatch []string, year int, variables, classification, cropCode string) string {
	ufParam := strings.Join(ufBatch, ",")
	path := fmt.Sprintf(
		"/t/%s/n6/in n3 %s/p/%d/v/%s/c%s/%s",
		table,
		ufParam,
		year,
		variables,
		classification,
		cropCode,
	)
	return sidraValuesBase + strings.ReplaceAll(path, " ", "%20")
}

func parseSIDRARows(raw []byte) ([]map[string]any, error) {
	var payload []map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	if len(payload) <= 1 {
		return nil, nil
	}
	return payload[1:], nil
}

func resolveCropCodes(entry catalog.RegistryEntry, crop string) (map[string]string, error) {
	if len(entry.SidraCrops) == 0 {
		return nil, fmt.Errorf("dataset %s missing sidra_crops", entry.DatasetID)
	}

	crop = strings.ToLower(strings.TrimSpace(crop))
	if crop == "" || crop == "all" {
		out := make(map[string]string, len(entry.SidraCrops))
		for name, code := range entry.SidraCrops {
			out[strings.ToLower(name)] = strconv.Itoa(code)
		}
		return out, nil
	}

	code, ok := entry.SidraCrops[crop]
	if !ok {
		return nil, fmt.Errorf("unknown crop %q for %s (allowed: %v)", crop, entry.DatasetID, cropNames(entry.SidraCrops))
	}
	return map[string]string{crop: strconv.Itoa(code)}, nil
}

func resolveYearRange(entry catalog.RegistryEntry, from, to int) (int, int, error) {
	start := entry.PeriodStart
	if start == 0 {
		start = 2010
	}
	end := entry.PeriodEnd
	if end == 0 {
		end = time.Now().UTC().Year() - 1
	}
	if from > 0 {
		start = from
	}
	if to > 0 {
		end = to
	}
	if start > end {
		return 0, 0, fmt.Errorf("invalid year range %d-%d", start, end)
	}
	return start, end, nil
}

func chunkUFs(requested []string) [][]string {
	if len(requested) == 0 {
		return defaultUFChunks
	}
	return [][]string{requested}
}

func formatVariables(vars []int) string {
	if len(vars) == 0 {
		return ""
	}
	parts := make([]string, len(vars))
	for i, v := range vars {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ",")
}

func cropNames(crops map[string]int) []string {
	names := make([]string, 0, len(crops))
	for name := range crops {
		names = append(names, name)
	}
	return names
}

// ResolvePAMURL validates the catalog base URL for a PAM dataset.
func ResolvePAMURL(entry catalog.RegistryEntry) (string, error) {
	raw := strings.TrimSpace(entry.SourceURL)
	if raw == "" {
		return "", fmt.Errorf("dataset %s has no source_url", entry.DatasetID)
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid source_url for %s: %w", entry.DatasetID, err)
	}
	if parsed.Scheme != "https" {
		return "", fmt.Errorf("source_url for %s must use https", entry.DatasetID)
	}
	if !strings.EqualFold(parsed.Host, "apisidra.ibge.gov.br") {
		return "", fmt.Errorf("source_url for %s must be on apisidra.ibge.gov.br", entry.DatasetID)
	}
	return parsed.String(), nil
}

func isPAMDataset(datasetID string) bool {
	return strings.HasPrefix(datasetID, "ibge.pam-")
}
