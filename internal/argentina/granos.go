package argentina

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultSeriesAPIBase = "https://apis.datos.gob.ar/series/api/series/"

var defaultGranosSeries = []string{
	"AGRO_A_Soja_0003",
	"AGRO_A_Maiz_0003",
	"AGRO_A_Trigo_0003",
}

var granosSlug = map[string]string{
	"AGRO_A_Soja_0003":  "soja",
	"AGRO_A_Maiz_0003":  "milho",
	"AGRO_A_Trigo_0003": "trigo",
}

type granosRow struct {
	SeriesID    string `json:"series_id"`
	CommoditySlug string `json:"commodity_slug"`
	RefYear     string `json:"refyear"`
	Value       string `json:"value"`
	Unit        string `json:"unit"`
	Source      string `json:"source"`
}

type seriesAPIResponse struct {
	Data [][]any `json:"data"`
	Meta []struct {
		Field struct {
			ID          string `json:"id"`
			Description string `json:"description"`
			Units       string `json:"units"`
		} `json:"field"`
		Dataset struct {
			Source string `json:"source"`
		} `json:"dataset"`
	} `json:"meta"`
}

// FetchGranosSnapshot downloads MAGyP grain production series from datos.gob.ar.
func (c *Client) FetchGranosSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	seriesIDs := entry.ArgentinaSeriesIDs
	if len(seriesIDs) == 0 {
		seriesIDs = defaultGranosSeries
	}

	query := url.Values{}
	query.Set("ids", strings.Join(seriesIDs, ","))
	if strings.TrimSpace(fromDate) != "" {
		query.Set("start_date", strings.TrimSpace(fromDate))
	} else if entry.PeriodStart > 0 {
		query.Set("start_date", fmt.Sprintf("%04d-01-01", entry.PeriodStart))
	}

	sourceURL := defaultSeriesAPIBase + "?" + query.Encode()
	raw, err := c.download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	rows, err := parseGranosJSON(raw, seriesIDs)
	if err != nil {
		return nil, "", err
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("argentina granos returned no rows for %s", entry.DatasetID)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].CommoditySlug != rows[j].CommoditySlug {
			return rows[i].CommoditySlug < rows[j].CommoditySlug
		}
		return rows[i].RefYear < rows[j].RefYear
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	return payload, sourceURL, nil
}

func parseGranosJSON(raw []byte, seriesIDs []string) ([]granosRow, error) {
	var payload seriesAPIResponse
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("parse argentina granos json: %w", err)
	}
	if len(payload.Data) == 0 {
		return nil, fmt.Errorf("empty argentina series data")
	}

	metaByIndex := make([]struct {
		id     string
		unit   string
		source string
	}, len(payload.Meta))
	for i, meta := range payload.Meta {
		metaByIndex[i] = struct {
			id     string
			unit   string
			source string
		}{
			id:     strings.TrimSpace(meta.Field.ID),
			unit:   strings.TrimSpace(meta.Field.Units),
			source: strings.TrimSpace(meta.Dataset.Source),
		}
	}

	var rows []granosRow
	for _, point := range payload.Data {
		if len(point) < 2 {
			continue
		}
		refYear, ok := point[0].(string)
		if !ok {
			continue
		}
		year := refYear
		if len(year) >= 4 {
			year = year[:4]
		}

		for col := 1; col < len(point); col++ {
			metaIdx := col - 1
			if metaIdx >= len(metaByIndex) {
				break
			}
			value, ok := numericString(point[col])
			if !ok {
				continue
			}
			seriesID := metaByIndex[metaIdx].id
			if seriesID == "" && metaIdx < len(seriesIDs) {
				seriesID = seriesIDs[metaIdx]
			}
			rows = append(rows, granosRow{
				SeriesID:      seriesID,
				CommoditySlug: commoditySlug(seriesID),
				RefYear:       year,
				Value:         value,
				Unit:          metaByIndex[metaIdx].unit,
				Source:        metaByIndex[metaIdx].source,
			})
		}
	}
	return rows, nil
}

func numericString(raw any) (string, bool) {
	switch v := raw.(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case int:
		return strconv.Itoa(v), true
	case json.Number:
		return v.String(), true
	default:
		return "", false
	}
}

func commoditySlug(seriesID string) string {
	if slug, ok := granosSlug[seriesID]; ok {
		return slug
	}
	return strings.ToLower(strings.TrimSpace(seriesID))
}

// FlattenGranos converts merged MAGyP grain production JSON into bronze columns.
func FlattenGranos(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []granosRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse argentina granos json: %w", err)
	}

	headers := []string{"series_id", "commodity_slug", "refyear", "value", "unit", "source"}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.RefYear) == "" || strings.TrimSpace(row.Value) == "" {
			continue
		}
		out = append(out, []string{
			row.SeriesID,
			row.CommoditySlug,
			row.RefYear,
			row.Value,
			row.Unit,
			row.Source,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no argentina granos rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

// ResolveGranosURL returns the datos.gob.ar series API URL for grain production.
func ResolveGranosURL(entry catalog.RegistryEntry) (string, error) {
	seriesIDs := entry.ArgentinaSeriesIDs
	if len(seriesIDs) == 0 {
		seriesIDs = defaultGranosSeries
	}
	query := url.Values{}
	query.Set("ids", strings.Join(seriesIDs, ","))
	if entry.PeriodStart > 0 {
		query.Set("start_date", fmt.Sprintf("%04d-01-01", entry.PeriodStart))
	}
	return defaultSeriesAPIBase + "?" + query.Encode(), nil
}
