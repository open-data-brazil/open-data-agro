package fao

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

type giewsRow struct {
	CountryCode      string `json:"country_code"`
	CountryName      string `json:"country_name"`
	CropSlug         string `json:"crop_slug"`
	MarketingYear    string `json:"marketing_year"`
	ProductionTrend  string `json:"production_trend"`
	OutlookNote      string `json:"outlook_note"`
}

// FetchGIEWSCropProspectsSnapshot returns FAO GIEWS crop prospects (fixture — FPMA CSV requires session).
func (c *Client) FetchGIEWSCropProspectsSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	_ = ctx
	rows, err := embeddedGIEWSample()
	if err != nil {
		return nil, "", err
	}
	if path := strings.TrimSpace(os.Getenv("GIEWS_BULK_PATH")); path != "" {
		if parsed, parseErr := parseGIEWSJSONFile(path); parseErr == nil && len(parsed) > 0 {
			rows = parsed
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].CountryCode+rows[i].CropSlug < rows[j].CountryCode+rows[j].CropSlug
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = "https://www.fao.org/giews/food-prospects-archive/en/"
	}
	return payload, sourceURL + " (fixture — FPMA bulk CSV requires browser session 2026-06-26)", nil
}

func embeddedGIEWSample() ([]giewsRow, error) {
	return []giewsRow{
		{CountryCode: "BRA", CountryName: "Brazil", CropSlug: "milho", MarketingYear: "2024", ProductionTrend: "stable", OutlookNote: "Favorable rains in Center-West"},
		{CountryCode: "ARG", CountryName: "Argentina", CropSlug: "soja", MarketingYear: "2024", ProductionTrend: "down", OutlookNote: "Dry conditions in Pampas"},
	}, nil
}

func parseGIEWSJSONFile(path string) ([]giewsRow, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rows []giewsRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// FlattenGIEWSCropProspects converts merged GIEWS JSON into canonical bronze columns.
func FlattenGIEWSCropProspects(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []giewsRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse fao giews json: %w", err)
	}

	headers := []string{"country_code", "country_name", "crop_slug", "marketing_year", "production_trend", "outlook_note"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.CountryCode) == "" {
			continue
		}
		out = append(out, []string{
			row.CountryCode, row.CountryName, row.CropSlug, row.MarketingYear,
			row.ProductionTrend, row.OutlookNote,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no fao giews rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
