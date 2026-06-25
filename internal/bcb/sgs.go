package bcb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const maxChunkYears = 10

type sgsRow struct {
	Data  string `json:"data"`
	Valor string `json:"valor"`
}

// FetchSGSSnapshot downloads and merges paginated SGS JSON observations.
func (c *Client) FetchSGSSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	code := entry.SGSCode
	if code == 0 {
		return nil, "", fmt.Errorf("dataset %s missing sgs_code", entry.DatasetID)
	}

	start, end, err := resolveDateRange(entry)
	if err != nil {
		return nil, "", err
	}

	chunks := chunkDateRange(start, end, maxChunkYears)
	merged := make(map[string]sgsRow)
	var requestURLs []string

	for _, chunk := range chunks {
		requestURL := buildSGSURL(code, chunk.from, chunk.to)
		requestURLs = append(requestURLs, requestURL)

		result, err := c.Download(ctx, requestURL)
		if err != nil {
			return nil, "", fmt.Errorf("sgs fetch %s-%s: %w", chunk.from, chunk.to, err)
		}

		var rows []sgsRow
		if err := json.Unmarshal(result.Body, &rows); err != nil {
			return nil, "", fmt.Errorf("parse sgs response %s: %w", requestURL, err)
		}
		for _, row := range rows {
			merged[row.Data] = row
		}
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("sgs returned no data rows for %s", entry.DatasetID)
	}

	keys := make([]string, 0, len(merged))
	for key := range merged {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	out := make([]sgsRow, 0, len(keys))
	for _, key := range keys {
		out = append(out, merged[key])
	}

	payload, err := json.Marshal(out)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s.%d/dados (chunks: %d)", sgsAPIBase, code, len(requestURLs))
	return payload, sourceURL, nil
}

type dateChunk struct {
	from string
	to   string
}

func buildSGSURL(code int, dataInicial, dataFinal string) string {
	values := url.Values{}
	values.Set("formato", "json")
	values.Set("dataInicial", dataInicial)
	values.Set("dataFinal", dataFinal)
	return fmt.Sprintf("%s.%d/dados?%s", sgsAPIBase, code, values.Encode())
}

func resolveDateRange(entry catalog.RegistryEntry) (time.Time, time.Time, error) {
	startRaw := strings.TrimSpace(entry.StartDate)
	if startRaw == "" && entry.PeriodStart > 0 {
		startRaw = fmt.Sprintf("01/01/%d", entry.PeriodStart)
	}
	if startRaw == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("dataset %s missing start_date or period_start", entry.DatasetID)
	}

	start, err := parseBCBDate(startRaw)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date for %s: %w", entry.DatasetID, err)
	}

	end := time.Now().UTC()
	if entry.PeriodEnd > 0 {
		candidate := time.Date(entry.PeriodEnd, 12, 31, 0, 0, 0, 0, time.UTC)
		if candidate.Before(end) {
			end = candidate
		}
	}
	if start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date range for %s", entry.DatasetID)
	}
	return start, end, nil
}

func chunkDateRange(start, end time.Time, maxYears int) []dateChunk {
	var chunks []dateChunk
	cursor := start
	for !cursor.After(end) {
		chunkEnd := cursor.AddDate(maxYears, 0, -1)
		if chunkEnd.After(end) {
			chunkEnd = end
		}
		chunks = append(chunks, dateChunk{
			from: formatBCBDate(cursor),
			to:   formatBCBDate(chunkEnd),
		})
		cursor = chunkEnd.AddDate(0, 0, 1)
	}
	return chunks
}

func parseBCBDate(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
 layouts := []string{"02/01/2006", "2006-01-02"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unparseable date %q", raw)
}

func formatBCBDate(t time.Time) string {
	return t.Format("02/01/2006")
}

// FlattenSGS converts merged SGS JSON rows into canonical bronze columns.
func FlattenSGS(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []sgsRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse sgs json: %w", err)
	}

	code := strconv.Itoa(entry.SGSCode)
	headers := []string{"sgs_codigo", "data", "valor", "ano"}
	out := make([][]string, 0, len(rows))

	for _, row := range rows {
		isoDate, err := normalizeObservationDate(row.Data)
		if err != nil {
			continue
		}
		valor, err := normalizeValor(row.Valor)
		if err != nil {
			continue
		}
		out = append(out, []string{
			code,
			isoDate,
			valor,
			isoDate[:4],
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no sgs rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

func normalizeObservationDate(raw string) (string, error) {
	t, err := parseBCBDate(raw)
	if err != nil {
		return "", err
	}
	return t.Format("2006-01-02"), nil
}

func normalizeValor(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty valor")
	}
	return strings.ReplaceAll(raw, ",", "."), nil
}
