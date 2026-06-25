package inmet

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// BDMEPFetchOptions controls annual ZIP pulls and UF filtering.
type BDMEPFetchOptions struct {
	Year int
	UFs  []string
}

// FetchBDMEPDailySnapshot downloads an annual ZIP and returns long-format daily CSV bytes.
func (c *Client) FetchBDMEPDailySnapshot(ctx context.Context, entry catalog.RegistryEntry, opts BDMEPFetchOptions) ([]byte, string, error) {
	sourceURL, err := ResolveAnnualZIPURL(entry, opts.Year)
	if err != nil {
		return nil, "", err
	}

	result, err := c.Download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	allowedVars := allowedClimateVariables(entry)
	ufFilter := normalizeUFs(opts.UFs, entry.PriorityUFs)

	rows, err := parseAnnualZIPToDailyLong(result.Body, opts.Year, ufFilter, allowedVars)
	if err != nil {
		return nil, "", err
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("no daily rows parsed for %s year=%d ufs=%v", entry.DatasetID, opts.Year, ufFilter)
	}

	payload, err := writeLongCSV(rows)
	if err != nil {
		return nil, "", err
	}
	return payload, sourceURL, nil
}

// FetchBDMEPMonthlySnapshot aggregates daily long rows to month buckets.
func (c *Client) FetchBDMEPMonthlySnapshot(ctx context.Context, entry catalog.RegistryEntry, opts BDMEPFetchOptions) ([]byte, string, error) {
	daily, sourceURL, err := c.FetchBDMEPDailySnapshot(ctx, entry, opts)
	if err != nil {
		return nil, "", err
	}
	records, err := readLatinCSV(daily)
	if err != nil {
		return nil, "", err
	}
	if len(records) <= 1 {
		return nil, "", fmt.Errorf("no daily rows to aggregate for %s", entry.DatasetID)
	}

	type monthKey struct {
		station  string
		month    string
		variable string
		uf       string
		year     string
	}
	accumulators := make(map[monthKey]*dailyAccumulator)

	for _, record := range records[1:] {
		if len(record) < 6 {
			continue
		}
		month := record[1]
		if len(month) >= 7 {
			month = month[:7]
		}
		key := monthKey{
			station:  record[0],
			month:    month,
			variable: record[2],
			uf:       record[4],
			year:     record[5],
		}
		value, ok := parseNumeric(record[3])
		if !ok {
			continue
		}
		acc, exists := accumulators[key]
		if !exists {
			agg := "mean"
			if key.variable == "precipitacao" {
				agg = "sum"
			}
			acc = &dailyAccumulator{agg: agg}
			accumulators[key] = acc
		}
		acc.add(value, acc.agg)
	}

	rows := make([][]string, 0, len(accumulators))
	for key, acc := range accumulators {
		value, ok := acc.value()
		if !ok {
			continue
		}
		rows = append(rows, []string{key.station, key.month, key.variable, value, key.uf, key.year})
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("no monthly rows aggregated for %s", entry.DatasetID)
	}

	payload, err := writeMonthlyCSV(rows)
	if err != nil {
		return nil, "", err
	}
	return payload, sourceURL, nil
}

func writeMonthlyCSV(rows [][]string) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"cd_estacao", "mes", "variavel", "valor", "uf", "ano"}); err != nil {
		return nil, err
	}
	if err := writer.WriteAll(rows); err != nil {
		return nil, err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func parseAnnualZIPToDailyLong(zipBytes []byte, year int, ufFilter map[string]struct{}, allowedVars map[string]struct{}) ([][]string, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, fmt.Errorf("open annual zip: %w", err)
	}

	var rows [][]string
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToUpper(file.Name), ".CSV") {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		raw, err := io.ReadAll(io.LimitReader(rc, 32<<20))
		_ = rc.Close()
		if err != nil {
			return nil, err
		}

		meta, dataRaw, err := splitStationFile(raw)
		if err != nil {
			continue
		}
		if len(ufFilter) > 0 {
			if _, ok := ufFilter[strings.ToUpper(meta.State)]; !ok {
				continue
			}
		}
		stationRows, err := parseStationDailyLong(meta, dataRaw, year, allowedVars)
		if err != nil || len(stationRows) == 0 {
			continue
		}
		rows = append(rows, stationRows...)
	}
	return rows, nil
}

func allowedClimateVariables(entry catalog.RegistryEntry) map[string]struct{} {
	if len(entry.ClimateVariables) == 0 {
		return nil
	}
	out := make(map[string]struct{}, len(entry.ClimateVariables))
	for _, name := range entry.ClimateVariables {
		out[strings.ToLower(strings.TrimSpace(name))] = struct{}{}
	}
	return out
}

func normalizeUFs(requested, defaults []string) map[string]struct{} {
	source := requested
	if len(source) == 0 {
		source = defaults
	}
	if len(source) == 0 {
		return nil
	}
	out := make(map[string]struct{}, len(source))
	for _, uf := range source {
		uf = strings.ToUpper(strings.TrimSpace(uf))
		if uf != "" {
			out[uf] = struct{}{}
		}
	}
	return out
}

func isBDMEPDataset(datasetID string) bool {
	switch datasetID {
	case "inmet.bdmep-diario", "inmet.bdmep-mensal", "inmet.pacote-anual-automaticas":
		return true
	default:
		return false
	}
}

func isStationCatalogDataset(datasetID string) bool {
	return datasetID == "inmet.estacoes-automaticas" || datasetID == "inmet.estacoes-convencionais"
}

// FlattenINMETCSV routes INMET CSV payloads to the correct flattener.
func FlattenINMETCSV(datasetID string, raw []byte) ([]string, [][]string, error) {
	switch datasetID {
	case "inmet.estacoes-automaticas", "inmet.estacoes-convencionais":
		return FlattenEstacoes(datasetID, raw)
	case "inmet.bdmep-diario", "inmet.pacote-anual-automaticas":
		return flattenLongCSV(raw, []string{"cd_estacao", "data", "variavel", "valor", "uf", "ano"})
	case "inmet.bdmep-mensal":
		return flattenLongCSV(raw, []string{"cd_estacao", "mes", "variavel", "valor", "uf", "ano"})
	default:
		return nil, nil, fmt.Errorf("unsupported inmet dataset %s", datasetID)
	}
}

func flattenLongCSV(raw []byte, expected []string) ([]string, [][]string, error) {
	records, err := readLatinCSV(raw)
	if err != nil {
		return nil, nil, err
	}
	if len(records) <= 1 {
		return nil, nil, fmt.Errorf("no data rows in long csv")
	}
	width := len(expected)
	rows := make([][]string, 0, len(records)-1)
	for _, record := range records[1:] {
		row := make([]string, width)
		copy(row, record)
		rows = append(rows, row)
	}
	return expected, rows, nil
}
