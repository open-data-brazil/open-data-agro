package b3

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

const sprdDownloadBase = "https://www.b3.com.br/pesquisapregao/download?filelist="

type futuroRow struct {
	RefDate        string `json:"refdate"`
	Symbol         string `json:"symbol"`
	Commodity      string `json:"commodity"`
	MaturityCode   string `json:"maturity_code"`
	PreviousPrice  string `json:"previous_price"`
	Price          string `json:"price"`
	Currency       string `json:"currency"`
	pendingDate    bool
	pendingSymbol  bool
	pendingAdj     bool
	pendingPrev    bool
}

// FetchFuturoSnapshot downloads daily settlement rows from B3 SPRD files.
func (c *Client) FetchFuturoSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	prefix := strings.ToUpper(strings.TrimSpace(entry.B3CommodityPrefix))
	if prefix == "" {
		return nil, "", fmt.Errorf("dataset %s missing b3_commodity_prefix", entry.DatasetID)
	}
	filePrefix := strings.ToUpper(strings.TrimSpace(entry.B3FilePrefix))
	if filePrefix == "" {
		filePrefix = "SPRD"
	}

	start, end, err := resolveFuturoRange(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	merged := make(map[string]futuroRow)
	downloads := 0

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		fileName := fmt.Sprintf("%s%s.zip", filePrefix, d.Format("060102"))
		url := sprdDownloadBase + fileName

		result, err := c.Download(ctx, url)
		if err != nil {
			continue
		}

		rows, err := parseSPRDZip(result.Body, prefix)
		if err != nil {
			continue
		}
		for _, row := range rows {
			key := row.RefDate + "|" + row.Symbol
			merged[key] = row
		}
		downloads++
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("b3 sprd returned no %s rows for %s", prefix, entry.DatasetID)
	}

	keys := make([]string, 0, len(merged))
	for key := range merged {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	out := make([]futuroRow, 0, len(keys))
	for _, key := range keys {
		out = append(out, merged[key])
	}

	payload, err := json.Marshal(out)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s%s (commodity=%s, days=%d)", sprdDownloadBase, filePrefix, prefix, downloads)
	return payload, sourceURL, nil
}

func resolveFuturoRange(entry catalog.RegistryEntry, fromDate string) (time.Time, time.Time, error) {
	startYear := entry.PeriodStart
	if startYear == 0 {
		startYear = 2024
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

// FlattenFuturo converts merged B3 futures JSON into canonical bronze columns.
func FlattenFuturo(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []futuroRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse b3 futuro json: %w", err)
	}

	headers := []string{
		"refdate", "symbol", "commodity", "maturity_code",
		"previous_price", "price", "currency", "price_change",
	}
	out := make([][]string, 0, len(rows))

	for _, row := range rows {
		if strings.TrimSpace(row.RefDate) == "" || strings.TrimSpace(row.Symbol) == "" {
			continue
		}
		change, err := priceChange(row.PreviousPrice, row.Price)
		if err != nil {
			change = ""
		}
		out = append(out, []string{
			row.RefDate,
			row.Symbol,
			row.Commodity,
			row.MaturityCode,
			row.PreviousPrice,
			row.Price,
			row.Currency,
			change,
		})
	}

	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no b3 futuro rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

func priceChange(previous, current string) (string, error) {
	prev, err := strconv.ParseFloat(strings.TrimSpace(previous), 64)
	if err != nil {
		return "", err
	}
	cur, err := strconv.ParseFloat(strings.TrimSpace(current), 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(cur-prev, 'f', -1, 64), nil
}
