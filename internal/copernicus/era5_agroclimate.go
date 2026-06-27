package copernicus

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

type era5Row struct {
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	RefDate      string `json:"refdate"`
	VariableSlug string `json:"variable_slug"`
	Value        string `json:"value"`
	Unit         string `json:"unit"`
}

// Client accesses Copernicus Climate Data Store (CDS) ERA5 agroclimate extracts.
type Client struct{}

// NewClient creates a Copernicus CDS client stub.
func NewClient() *Client { return &Client{} }

// FetchERA5AgroclimateSnapshot returns ERA5 point sample (API stub until CDS key configured).
func (c *Client) FetchERA5AgroclimateSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	_ = ctx
	if strings.TrimSpace(os.Getenv("COPERNICUS_CDS_API_KEY")) != "" {
		return nil, "", fmt.Errorf("copernicus cds live fetch not implemented — use fixture stub (Phase 49 P2)")
	}

	rows, err := embeddedERA5Sample()
	if err != nil {
		return nil, "", err
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].RefDate+rows[i].VariableSlug < rows[j].RefDate+rows[j].VariableSlug
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = "https://cds.climate.copernicus.eu/"
	}
	return payload, sourceURL + " (fixture stub — set COPERNICUS_CDS_API_KEY for live CDS)", nil
}

func embeddedERA5Sample() ([]era5Row, error) {
	return []era5Row{
		{Latitude: "-15.8", Longitude: "-47.9", RefDate: "2024-01-01", VariableSlug: "t2m", Value: "24.0", Unit: "degC"},
		{Latitude: "-15.8", Longitude: "-47.9", RefDate: "2024-01-01", VariableSlug: "tp", Value: "22.44", Unit: "mm"},
	}, nil
}

// FlattenERA5Agroclimate converts merged ERA5 JSON into canonical bronze columns.
func FlattenERA5Agroclimate(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []era5Row
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse copernicus era5 json: %w", err)
	}

	headers := []string{"latitude", "longitude", "refdate", "variable_slug", "value", "unit"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RefDate) == "" {
			continue
		}
		out = append(out, []string{
			row.Latitude, row.Longitude, row.RefDate, row.VariableSlug, row.Value, row.Unit,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no copernicus era5 rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
