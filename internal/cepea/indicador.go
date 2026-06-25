package cepea

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// FetchIndicadorSnapshot downloads and parses a CEPEA daily indicator series.
func (c *Client) FetchIndicadorSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	spec, err := IndicatorSpecFor(entry.DatasetID.String())
	if err != nil {
		return nil, "", err
	}

	start, err := resolveFromDate(entry, fromDate)
	if err != nil {
		return nil, "", err
	}

	portalURL, err := PortalIndicatorURL(entry.DatasetID.String())
	if err != nil {
		return nil, "", err
	}

	mirrorURL, err := MirrorURL(entry.DatasetID.String())
	if err != nil {
		return nil, "", err
	}

	var observations []Observation
	var sourceURL string
	var lastModified string

	result, err := c.Download(ctx, portalURL)
	if err == nil {
		observations, err = ParseIndicatorHTML(result.Body, spec.Praca)
		if err == nil {
			sourceURL = portalURL
			lastModified = result.LastModified
		}
	}

	if len(observations) == 0 {
		mirrorResult, mirrorErr := c.Download(ctx, mirrorURL)
		if mirrorErr != nil {
			if err != nil {
				return nil, "", fmt.Errorf("cepea portal: %v; mirror: %w", err, mirrorErr)
			}
			return nil, "", mirrorErr
		}
		observations, err = ParseIndicatorHTML(mirrorResult.Body, spec.Praca)
		if err != nil {
			return nil, "", fmt.Errorf("parse mirror html: %w", err)
		}
		sourceURL = mirrorURL + " (CEPEA mirror; origin CEPEA/ESALQ)"
		lastModified = mirrorResult.LastModified
	}

	filtered := filterFromDate(observations, start)
	if len(filtered) == 0 {
		return nil, "", fmt.Errorf("no observations on or after %s for %s", start.Format("2006-01-02"), entry.DatasetID)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Data < filtered[j].Data
	})

	payload, err := json.Marshal(filtered)
	if err != nil {
		return nil, "", err
	}

	if lastModified != "" {
		sourceURL = fmt.Sprintf("%s [Last-Modified: %s]", sourceURL, lastModified)
	}
	return payload, sourceURL, nil
}

func resolveFromDate(entry catalog.RegistryEntry, override string) (time.Time, error) {
	raw := strings.TrimSpace(override)
	if raw == "" {
		raw = strings.TrimSpace(entry.StartDate)
	}
	if raw == "" && entry.PeriodStart > 0 {
		raw = fmt.Sprintf("%d-01-01", entry.PeriodStart)
	}
	if raw == "" {
		return time.Time{}, fmt.Errorf("dataset %s missing start_date", entry.DatasetID)
	}

	layouts := []string{"2006-01-02", "02/01/2006"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid from date %q for %s", raw, entry.DatasetID)
}

func filterFromDate(rows []Observation, start time.Time) []Observation {
	startKey := start.Format("2006-01-02")
	out := make([]Observation, 0, len(rows))
	for _, row := range rows {
		if row.Data >= startKey {
			out = append(out, row)
		}
	}
	return out
}

// FlattenIndicador converts parsed JSON observations into canonical bronze columns.
func FlattenIndicador(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []Observation
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse cepea json: %w", err)
	}

	spec, err := IndicatorSpecFor(entry.DatasetID.String())
	if err != nil {
		return nil, nil, err
	}

	headers := []string{
		"produto",
		"praca",
		"data",
		"preco_rs_sc",
		"variacao_dia_pct",
		"preco_usd_sc",
		"ano",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if row.Data == "" || row.PrecoRsSc == "" {
			continue
		}
		out = append(out, []string{
			spec.ProductSlug,
			spec.Praca,
			row.Data,
			row.PrecoRsSc,
			row.VariacaoDiaPct,
			row.PrecoUsdSc,
			row.Data[:4],
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no cepea rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
