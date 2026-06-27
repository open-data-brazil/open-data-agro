package noaa

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

type gpccRow struct {
	RefMonth        string `json:"refmonth"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
	PrecipMm        string `json:"precip_mm"`
	GridResolution  string `json:"grid_resolution"`
}

// FetchGPCCPrecipitationSnapshot returns GPCC grid precipitation sample (fixture — bulk URL unstable).
func (c *Client) FetchGPCCPrecipitationSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	_ = ctx
	rows, err := embeddedGPCCSample()
	if err != nil {
		return nil, "", err
	}
	if path := strings.TrimSpace(os.Getenv("GPCC_BULK_PATH")); path != "" {
		if parsed, parseErr := parseGPCCJSONFile(path); parseErr == nil && len(parsed) > 0 {
			rows = parsed
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].RefMonth < rows[j].RefMonth
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	sourceURL := strings.TrimSpace(entry.NOAAIndexURL)
	if sourceURL == "" {
		sourceURL = "https://opendata.dwd.de/climate_environment/GPCC/Monitoring/"
	}
	return payload, sourceURL + " (fixture — GPCC monthly bulk URL unstable 2026-06-26)", nil
}

func embeddedGPCCSample() ([]gpccRow, error) {
	return []gpccRow{
		{RefMonth: "2024-01-01", Latitude: "-15.75", Longitude: "-47.75", PrecipMm: "145.2", GridResolution: "1.0"},
		{RefMonth: "2024-02-01", Latitude: "-15.75", Longitude: "-47.75", PrecipMm: "132.8", GridResolution: "1.0"},
	}, nil
}

func parseGPCCJSONFile(path string) ([]gpccRow, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rows []gpccRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// FlattenGPCCPrecipitation converts merged GPCC JSON into canonical bronze columns.
func FlattenGPCCPrecipitation(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []gpccRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse noaa gpcc json: %w", err)
	}

	headers := []string{"refmonth", "latitude", "longitude", "precip_mm", "grid_resolution"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RefMonth) == "" {
			continue
		}
		out = append(out, []string{
			row.RefMonth, row.Latitude, row.Longitude, row.PrecipMm, row.GridResolution,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no noaa gpcc rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
