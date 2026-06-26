package mdic

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const generalEndpoint = "/general"

type comexRequest struct {
	Flow        string        `json:"flow"`
	MonthDetail bool          `json:"monthDetail"`
	Period      comexPeriod   `json:"period"`
	Filters     []comexFilter `json:"filters"`
	Details     []string      `json:"details"`
	Metrics     []string      `json:"metrics"`
}

type comexPeriod struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type comexFilter struct {
	Filter string   `json:"filter"`
	Values []string `json:"values"`
}

type comexResponse struct {
	Success bool `json:"success"`
	Data    struct {
		List []comexRow `json:"list"`
	} `json:"data"`
}

type comexRow struct {
	CoNCM           string `json:"coNcm"`
	NCM             string `json:"ncm"`
	Year            string `json:"year"`
	MonthNumber     string `json:"monthNumber"`
	State           string `json:"state"`
	MetricFOB       string `json:"metricFOB"`
	MetricKG        string `json:"metricKG"`
	MetricCIF       string `json:"metricCIF"`
	MetricFreight   string `json:"metricFreight"`
	MetricInsurance string `json:"metricInsurance"`
}

type mergedRow struct {
	comexRow
	periodFrom string
	periodTo   string
}

// ProdutoSlug maps official NCM codes to canonical crop/product slugs.
var ProdutoSlug = map[string]string{
	"12019000": "soja",
	"10059000": "milho",
	"10019900": "trigo",
	"02013000": "carne_bovina",
	"02023000": "carne_bovina",
	"31021010": "ureia",
	"31022100": "sulfato_amonia",
	"31023000": "nitrato_amonia",
	"31024000": "misturas_amonia",
	"31052000": "npk",
	"27101921": "diesel",
	"27101922": "oleo_combustivel",
}

// StateToUF maps Comex Stat state names to IBGE sigla.
var StateToUF = map[string]string{
	"Acre":                "AC",
	"Alagoas":             "AL",
	"Amapá":               "AP",
	"Amazonas":            "AM",
	"Bahia":               "BA",
	"Ceará":               "CE",
	"Distrito Federal":    "DF",
	"Espírito Santo":      "ES",
	"Goiás":               "GO",
	"Maranhão":            "MA",
	"Mato Grosso":         "MT",
	"Mato Grosso do Sul":  "MS",
	"Minas Gerais":        "MG",
	"Pará":                "PA",
	"Paraíba":             "PB",
	"Paraná":              "PR",
	"Pernambuco":          "PE",
	"Piauí":               "PI",
	"Rio de Janeiro":      "RJ",
	"Rio Grande do Norte": "RN",
	"Rio Grande do Sul":   "RS",
	"Rondônia":            "RO",
	"Roraima":             "RR",
	"Santa Catarina":      "SC",
	"São Paulo":           "SP",
	"Sergipe":             "SE",
	"Tocantins":           "TO",
}

// FetchComexSnapshot downloads monthly Comex rows for configured NCM codes.
func (c *Client) FetchComexSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	flow := strings.TrimSpace(entry.ComexFlow)
	if flow == "" {
		flow = "export"
	}
	ncms := entry.ComexNCMs
	if len(ncms) == 0 {
		return nil, "", fmt.Errorf("dataset %s missing comex_ncms", entry.DatasetID)
	}

	details := entry.ComexDetails
	if len(details) == 0 {
		details = []string{"ncm"}
	}
	metrics := entry.ComexMetrics
	if len(metrics) == 0 {
		if strings.EqualFold(flow, "import") {
			metrics = []string{"metricCIF", "metricKG", "metricFreight", "metricInsurance"}
		} else {
			metrics = []string{"metricFOB", "metricKG"}
		}
	}

	start, end, err := resolveComexRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	hasState := false
	for _, d := range details {
		if strings.EqualFold(d, "state") {
			hasState = true
			break
		}
	}

	merged := make(map[string]mergedRow)
	var chunkCount int

	for year := start.Year(); year <= end.Year(); year++ {
		periodFrom := fmt.Sprintf("%04d-01", year)
		periodTo := fmt.Sprintf("%04d-12", year)
		if year == start.Year() {
			periodFrom = fmt.Sprintf("%04d-%02d", year, int(start.Month()))
		}
		if year == end.Year() {
			periodTo = fmt.Sprintf("%04d-%02d", year, int(end.Month()))
		}

		reqBody, err := json.Marshal(comexRequest{
			Flow:        flow,
			MonthDetail: true,
			Period:      comexPeriod{From: periodFrom, To: periodTo},
			Filters:     []comexFilter{{Filter: "ncm", Values: ncms}},
			Details:     details,
			Metrics:     metrics,
		})
		if err != nil {
			return nil, "", err
		}

		raw, err := c.PostJSON(ctx, generalEndpoint, reqBody)
		if err != nil {
			return nil, "", fmt.Errorf("comex fetch %s-%s: %w", periodFrom, periodTo, err)
		}

		var resp comexResponse
		if err := json.Unmarshal(raw, &resp); err != nil {
			return nil, "", fmt.Errorf("parse comex response %s-%s: %w", periodFrom, periodTo, err)
		}
		if !resp.Success && len(resp.Data.List) == 0 {
			continue
		}

		for _, row := range resp.Data.List {
			key := mergeKey(row, hasState)
			merged[key] = mergedRow{comexRow: row, periodFrom: periodFrom, periodTo: periodTo}
		}
		chunkCount++

		select {
		case <-ctx.Done():
			return nil, "", ctx.Err()
		case <-time.After(400 * time.Millisecond):
		}
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("comex returned no data rows for %s", entry.DatasetID)
	}

	keys := make([]string, 0, len(merged))
	for key := range merged {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	out := make([]comexRow, 0, len(keys))
	for _, key := range keys {
		out = append(out, merged[key].comexRow)
	}

	payload, err := json.Marshal(out)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s%s (flow=%s, ncm=%d, chunks=%d)", comexAPIBase, generalEndpoint, flow, len(ncms), chunkCount)
	return payload, sourceURL, nil
}

func mergeKey(row comexRow, hasState bool) string {
	if hasState {
		return row.CoNCM + "|" + strings.TrimSpace(row.State) + "|" + row.Year + "|" + row.MonthNumber
	}
	return row.CoNCM + "|" + row.Year + "|" + row.MonthNumber
}

func resolveComexRange(entry catalog.RegistryEntry, fromDate string) (time.Time, time.Time, error) {
	startYear := entry.PeriodStart
	if startYear == 0 {
		startYear = 2015
	}
	start := time.Date(startYear, 1, 1, 0, 0, 0, 0, time.UTC)

	if raw := strings.TrimSpace(fromDate); raw != "" {
		parsed, err := parseISODate(raw)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid from date for %s: %w", entry.DatasetID, err)
		}
		if parsed.After(start) {
			start = parsed
		}
	}

	end := time.Now().UTC()
	if entry.PeriodEnd > 0 {
		candidate := time.Date(entry.PeriodEnd, 12, 1, 0, 0, 0, 0, time.UTC)
		if candidate.Before(end) {
			end = candidate
		}
	}
	if start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date range for %s", entry.DatasetID)
	}
	return start, end, nil
}

func parseISODate(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	layouts := []string{"2006-01-02", "2006-01"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unparseable date %q", raw)
}

// FlattenComex converts merged Comex JSON rows into canonical bronze columns.
func FlattenComex(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	switch entry.DatasetID.String() {
	case "mdic.comex-importacao-ncm-mes", "mdic.comex-importacao-diesel-ncm":
		return flattenComexImport(entry, raw)
	case "mdic.comex-exportacao-uf-ncm":
		return flattenComexExportUF(entry, raw)
	default:
		return flattenComexExportNCM(entry, raw)
	}
}

func flattenComexExportNCM(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	rows, err := parseComexRows(raw)
	if err != nil {
		return nil, nil, err
	}

	headers := []string{
		"co_ncm", "ncm_descricao", "produto_slug", "data",
		"valor_fob_usd", "quantidade_kg", "ano",
	}
	out := make([][]string, 0, len(rows))

	for _, row := range rows {
		coNCM := strings.TrimSpace(row.CoNCM)
		if coNCM == "" {
			continue
		}
		isoDate, err := monthToDate(row.Year, row.MonthNumber)
		if err != nil {
			continue
		}
		fob, err := normalizeMetric(row.MetricFOB)
		if err != nil {
			continue
		}
		kg, err := normalizeMetric(row.MetricKG)
		if err != nil {
			continue
		}
		out = append(out, []string{
			coNCM,
			strings.TrimSpace(row.NCM),
			produtoSlug(coNCM),
			isoDate,
			fob,
			kg,
			isoDate[:4],
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no comex rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

func flattenComexImport(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	rows, err := parseComexRows(raw)
	if err != nil {
		return nil, nil, err
	}

	headers := []string{
		"co_ncm", "ncm_descricao", "produto_slug", "data",
		"valor_cif_usd", "quantidade_kg", "valor_frete_usd", "valor_seguro_usd", "ano",
	}
	out := make([][]string, 0, len(rows))

	for _, row := range rows {
		coNCM := strings.TrimSpace(row.CoNCM)
		if coNCM == "" {
			continue
		}
		isoDate, err := monthToDate(row.Year, row.MonthNumber)
		if err != nil {
			continue
		}
		cif, err := normalizeMetric(row.MetricCIF)
		if err != nil {
			continue
		}
		kg, err := normalizeMetric(row.MetricKG)
		if err != nil {
			continue
		}
		freight, err := normalizeMetric(row.MetricFreight)
		if err != nil {
			continue
		}
		insurance, err := normalizeMetric(row.MetricInsurance)
		if err != nil {
			continue
		}
		out = append(out, []string{
			coNCM,
			strings.TrimSpace(row.NCM),
			produtoSlug(coNCM),
			isoDate,
			cif,
			kg,
			freight,
			insurance,
			isoDate[:4],
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no comex rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

func flattenComexExportUF(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	rows, err := parseComexRows(raw)
	if err != nil {
		return nil, nil, err
	}

	headers := []string{
		"co_ncm", "ncm_descricao", "produto_slug", "uf", "data",
		"valor_fob_usd", "quantidade_kg", "ano",
	}
	out := make([][]string, 0, len(rows))

	for _, row := range rows {
		coNCM := strings.TrimSpace(row.CoNCM)
		if coNCM == "" {
			continue
		}
		uf := stateToUF(row.State)
		if uf == "" {
			continue
		}
		isoDate, err := monthToDate(row.Year, row.MonthNumber)
		if err != nil {
			continue
		}
		fob, err := normalizeMetric(row.MetricFOB)
		if err != nil {
			continue
		}
		kg, err := normalizeMetric(row.MetricKG)
		if err != nil {
			continue
		}
		out = append(out, []string{
			coNCM,
			strings.TrimSpace(row.NCM),
			produtoSlug(coNCM),
			uf,
			isoDate,
			fob,
			kg,
			isoDate[:4],
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no comex rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

func parseComexRows(raw []byte) ([]comexRow, error) {
	var rows []comexRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, fmt.Errorf("parse comex json: %w", err)
	}
	return rows, nil
}

func produtoSlug(coNCM string) string {
	if slug := ProdutoSlug[coNCM]; slug != "" {
		return slug
	}
	return "outros"
}

func stateToUF(name string) string {
	name = strings.TrimSpace(name)
	if uf, ok := StateToUF[name]; ok {
		return uf
	}
	return ""
}

func monthToDate(yearRaw, monthRaw string) (string, error) {
	year, err := strconv.Atoi(strings.TrimSpace(yearRaw))
	if err != nil || year < 1900 {
		return "", fmt.Errorf("invalid year %q", yearRaw)
	}
	month, err := strconv.Atoi(strings.TrimSpace(monthRaw))
	if err != nil || month < 1 || month > 12 {
		return "", fmt.Errorf("invalid month %q", monthRaw)
	}
	return fmt.Sprintf("%04d-%02d-01", year, month), nil
}

func normalizeMetric(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "0", nil
	}
	if _, err := strconv.ParseFloat(raw, 64); err != nil {
		return "", fmt.Errorf("invalid metric %q", raw)
	}
	return raw, nil
}
