package eurostat

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const defaultAgPricesDataset = "apri_pi15_outa"

type agPriceRow struct {
	DatasetCode string `json:"dataset_code"`
	Geo         string `json:"geo"`
	ProductCode string `json:"product_code"`
	ProductName string `json:"product_name"`
	Year        string `json:"year"`
	IndexValue  string `json:"index_value"`
	BasePeriod  string `json:"base_period"`
}

type eurostatDataset struct {
	ID        any                       `json:"id"`
	Value     map[string]float64        `json:"value"`
	Dimension map[string]eurostatDim    `json:"dimension"`
}

type eurostatDim struct {
	Category eurostatCategory `json:"category"`
}

type eurostatCategory struct {
	Index map[string]int            `json:"index"`
	Label map[string]string         `json:"label"`
}

// FetchAgPricesSnapshot downloads EU agricultural output price indices from Eurostat.
func (c *Client) FetchAgPricesSnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	datasetCode := strings.TrimSpace(entry.EurostatDatasetCode)
	if datasetCode == "" {
		datasetCode = defaultAgPricesDataset
	}
	geo := strings.TrimSpace(entry.EurostatGeo)
	if geo == "" {
		geo = "EU27_2020"
	}
	products := entry.EurostatProducts
	if len(products) == 0 {
		products = []string{"010000", "015000", "011000"}
	}
	sinceYear := entry.PeriodStart
	if sinceYear == 0 {
		sinceYear = 2020
	}
	if strings.TrimSpace(fromDate) != "" && len(fromDate) >= 4 {
		if y, err := strconv.Atoi(fromDate[:4]); err == nil && y > sinceYear {
			sinceYear = y
		}
	}

	sourceURL := buildDatasetURL(datasetCode, geo, products, sinceYear)
	raw, err := c.download(ctx, sourceURL)
	if err != nil {
		return nil, "", err
	}

	rows, err := parseAgPricesJSON(raw, datasetCode, geo)
	if err != nil {
		return nil, "", err
	}
	if len(rows) == 0 {
		return nil, "", fmt.Errorf("eurostat ag prices returned no rows for %s", entry.DatasetID)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].ProductCode != rows[j].ProductCode {
			return rows[i].ProductCode < rows[j].ProductCode
		}
		return rows[i].Year < rows[j].Year
	})

	payload, err := json.Marshal(rows)
	if err != nil {
		return nil, "", err
	}
	return payload, sourceURL, nil
}

func parseAgPricesJSON(raw []byte, datasetCode, geo string) ([]agPriceRow, error) {
	var payload eurostatDataset
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("parse eurostat json: %w", err)
	}
	if len(payload.Value) == 0 {
		return nil, fmt.Errorf("eurostat dataset has no values")
	}

	dimOrder := dimensionOrder(payload.ID)
	sizes := make([]int, len(dimOrder))
	labels := make([]map[string]string, len(dimOrder))
	indexes := make([]map[string]int, len(dimOrder))
	for i, name := range dimOrder {
		dim := payload.Dimension[name]
		sizes[i] = len(dim.Category.Index)
		labels[i] = dim.Category.Label
		indexes[i] = dim.Category.Index
	}

	var rows []agPriceRow
	for key, value := range payload.Value {
		idx, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		coords := decodeIndex(idx, sizes)
		dims := make(map[string]string)
		for i, name := range dimOrder {
			code := codeAt(indexes[i], coords[i])
			dims[name] = code
		}

		productCode := dims["product"]
		year := dims["time"]
		rowGeo := dims["geo"]
		if rowGeo == "" {
			rowGeo = geo
		}
		if productCode == "" || year == "" {
			continue
		}

		rows = append(rows, agPriceRow{
			DatasetCode: datasetCode,
			Geo:         rowGeo,
			ProductCode: productCode,
			ProductName: labelAt(labels, dimOrder, "product", productCode),
			Year:        year,
			IndexValue:  strconv.FormatFloat(value, 'f', -1, 64),
			BasePeriod:  "2015=100",
		})
	}
	return rows, nil
}

func dimensionOrder(id any) []string {
	switch v := id.(type) {
	case []any:
		out := make([]string, 0, len(v))
		for _, item := range v {
			out = append(out, fmt.Sprint(item))
		}
		return out
	case string:
		return strings.Split(v, ",")
	default:
		return []string{"freq", "p_adj", "unit", "product", "geo", "time"}
	}
}

func decodeIndex(flat int, sizes []int) []int {
	coords := make([]int, len(sizes))
	for i := len(sizes) - 1; i >= 0; i-- {
		size := sizes[i]
		if size == 0 {
			coords[i] = 0
			continue
		}
		coords[i] = flat % size
		flat /= size
	}
	return coords
}

func codeAt(index map[string]int, pos int) string {
	for code, idx := range index {
		if idx == pos {
			return code
		}
	}
	return ""
}

func dimIndex(order []string, name string) int {
	for i, item := range order {
		if item == name {
			return i
		}
	}
	return -1
}

func labelAt(all []map[string]string, order []string, dimName, code string) string {
	idx := dimIndex(order, dimName)
	if idx < 0 || idx >= len(all) {
		return ""
	}
	return all[idx][code]
}

// FlattenAgPrices converts merged Eurostat JSON into canonical bronze columns.
func FlattenAgPrices(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []agPriceRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse eurostat ag prices json: %w", err)
	}

	headers := []string{
		"dataset_code",
		"geo",
		"product_code",
		"product_name",
		"year",
		"index_value",
		"base_period",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.Year) == "" || strings.TrimSpace(row.IndexValue) == "" {
			continue
		}
		out = append(out, []string{
			row.DatasetCode,
			row.Geo,
			row.ProductCode,
			row.ProductName,
			row.Year,
			row.IndexValue,
			row.BasePeriod,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no eurostat ag prices rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}

// ResolveURL returns the Eurostat API URL for a catalog entry.
func ResolveURL(entry catalog.RegistryEntry) (string, error) {
	datasetCode := strings.TrimSpace(entry.EurostatDatasetCode)
	if datasetCode == "" {
		datasetCode = defaultAgPricesDataset
	}
	geo := strings.TrimSpace(entry.EurostatGeo)
	if geo == "" {
		geo = "EU27_2020"
	}
	return buildDatasetURL(datasetCode, geo, entry.EurostatProducts, entry.PeriodStart), nil
}
