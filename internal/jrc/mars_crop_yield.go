package jrc

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultMARSYieldURL = "https://agricultural-production-hotspots.ec.europa.eu/data/yield-forecast/recent/Central%20America_2024_75p.csv"

var cropSlug = map[string]string{
	"Maize (corn)": "milho",
	"Rice":         "arroz",
	"Wheat":        "trigo",
	"Soybeans":     "soja",
}

type marsRow struct {
	Country          string `json:"country"`
	Crop             string `json:"crop"`
	CropSlug         string `json:"crop_slug"`
	ForecastYield    string `json:"forecast_yield_kg_ha"`
	FiveYrAvg        string `json:"five_yr_avg_kg_ha"`
	HarvestYear      string `json:"harvest_year"`
	ForecastTiming   string `json:"forecast_timing"`
	RegionName       string `json:"region_name"`
}

// FetchMARSCropYieldSnapshot downloads JRC MARS ASAP yield forecast CSV.
func (c *Client) FetchMARSCropYieldSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = defaultMARSYieldURL
	}

	raw, err := c.download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	region, timing, year := parseMARSFilename(sourceURL)
	rows, err := parseMARSCSV(raw, region, timing, year)
	if err != nil {
		return nil, "", err
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("jrc mars returned no rows for %s", entry.DatasetID)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Country != rows[j].Country {
			return rows[i].Country < rows[j].Country
		}
		return rows[i].CropSlug < rows[j].CropSlug
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	return payload, sourceURL, nil
}

func parseMARSFilename(sourceURL string) (region, timing, year string) {
	base := sourceURL
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	base = strings.TrimSuffix(base, ".csv")
	parts := strings.Split(base, "_")
	if len(parts) >= 3 {
		region = parts[0]
		year = parts[1]
		timing = strings.TrimSuffix(parts[2], ".csv")
	}
	return region, timing, year
}

func parseMARSCSV(raw []byte, region, timing, year string) ([]marsRow, error) {
	reader := csv.NewReader(bytes.NewReader(raw))
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read jrc mars header: %w", err)
	}
	col := indexColumns(header)

	var rows []marsRow
	for {
		record, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, fmt.Errorf("read jrc mars row: %w", readErr)
		}
		country := field(record, col, "Country")
		crop := field(record, col, "Crop")
		forecast := field(record, col, "Forecasted yield (kg/ha)")
		avg := field(record, col, "5 yrs avg (kg/ha)")
		if country == "" || crop == "" || forecast == "" || forecast == "/" {
			continue
		}

		rows = append(rows, marsRow{
			Country:        country,
			Crop:           crop,
			CropSlug:       slugForCrop(crop),
			ForecastYield:  forecast,
			FiveYrAvg:      avg,
			HarvestYear:    year,
			ForecastTiming: timing,
			RegionName:     region,
		})
	}
	return rows, nil
}

func indexColumns(header []string) map[string]int {
	out := make(map[string]int, len(header))
	for i, name := range header {
		out[strings.TrimSpace(name)] = i
	}
	return out
}

func field(record []string, col map[string]int, name string) string {
	idx, ok := col[name]
	if !ok || idx >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[idx])
}

func slugForCrop(crop string) string {
	if slug, ok := cropSlug[crop]; ok {
		return slug
	}
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(crop), " ", "_"))
}

// FlattenMARSCropYield converts merged JRC MARS JSON into canonical bronze columns.
func FlattenMARSCropYield(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []marsRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse jrc mars json: %w", err)
	}

	headers := []string{
		"country", "crop", "crop_slug", "forecast_yield_kg_ha", "five_yr_avg_kg_ha",
		"harvest_year", "forecast_timing", "region_name",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.Country) == "" {
			continue
		}
		out = append(out, []string{
			row.Country, row.Crop, row.CropSlug, row.ForecastYield, row.FiveYrAvg,
			row.HarvestYear, row.ForecastTiming, row.RegionName,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no jrc mars rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
