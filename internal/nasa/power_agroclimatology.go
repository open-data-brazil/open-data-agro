package nasa

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultPOWERURL = "https://power.larc.nasa.gov/api/temporal/daily/point?parameters=PRECTOTCORR,T2M&community=AG&longitude=-47.9&latitude=-15.8&start=20240101&end=20240107&format=JSON"

type powerRow struct {
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	RefDate        string `json:"refdate"`
	ParameterSlug  string `json:"parameter_slug"`
	Value          string `json:"value"`
}

type powerResponse struct {
	Geometry struct {
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Parameter map[string]map[string]float64 `json:"parameter"`
	} `json:"properties"`
}

// FetchPOWERAgroSnapshot downloads NASA POWER daily agroclimatology point data.
func (c *Client) FetchPOWERAgroSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	sourceURL := strings.TrimSpace(entry.SourceURL)
	if sourceURL == "" {
		sourceURL = defaultPOWERURL
	}

	raw, err := c.download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	rows, err := parsePOWERJSON(raw)
	if err != nil {
		return nil, "", err
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("nasa power returned no rows for %s", entry.DatasetID)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].RefDate != rows[j].RefDate {
			return rows[i].RefDate < rows[j].RefDate
		}
		return rows[i].ParameterSlug < rows[j].ParameterSlug
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	return payload, sourceURL, nil
}

func parsePOWERJSON(raw []byte) ([]powerRow, error) {
	var resp powerResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse nasa power json: %w", err)
	}
	if len(resp.Geometry.Coordinates) < 2 {
		return nil, fmt.Errorf("nasa power missing coordinates")
	}
	lon := fmt.Sprintf("%.4f", resp.Geometry.Coordinates[0])
	lat := fmt.Sprintf("%.4f", resp.Geometry.Coordinates[1])

	var rows []powerRow
	for param, days := range resp.Properties.Parameter {
		slug := strings.ToLower(param)
		for day, val := range days {
			if len(day) != 8 {
				continue
			}
			refDate := fmt.Sprintf("%s-%s-%s", day[0:4], day[4:6], day[6:8])
			rows = append(rows, powerRow{
				Latitude:      lat,
				Longitude:     lon,
				RefDate:       refDate,
				ParameterSlug: slug,
				Value:         fmt.Sprintf("%.4f", val),
			})
		}
	}
	return rows, nil
}

// FlattenPOWERAgro converts merged NASA POWER JSON into canonical bronze columns.
func FlattenPOWERAgro(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []powerRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse nasa power json: %w", err)
	}

	headers := []string{"latitude", "longitude", "refdate", "parameter_slug", "value"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RefDate) == "" {
			continue
		}
		out = append(out, []string{
			row.Latitude, row.Longitude, row.RefDate, row.ParameterSlug, row.Value,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no nasa power rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
