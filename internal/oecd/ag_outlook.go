package oecd

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

var defaultCommodityCodes = []string{"CPC_0141", "CPC_0142", "CPC_0143"}

var defaultMeasureCodes = []string{"QP", "EX", "IM", "ST", "QC", "CR", "FE"}

type outlookRow struct {
	RefArea       string `json:"ref_area"`
	RefAreaName   string `json:"ref_area_name"`
	CommodityCode string `json:"commodity_code"`
	CommodityName string `json:"commodity_name"`
	MeasureCode   string `json:"measure_code"`
	MeasureName   string `json:"measure_name"`
	Unit          string `json:"unit"`
	UnitMult      string `json:"unit_mult"`
	Year          string `json:"year"`
	Value         string `json:"value"`
	ObsStatus     string `json:"obs_status"`
}

// FetchAgOutlookSnapshot downloads OECD-FAO Agricultural Outlook SDMX CSV for configured commodities.
func (c *Client) FetchAgOutlookSnapshot(ctx context.Context, entry catalog.RegistryEntry) ([]byte, string, error) {
	commodities := entry.OECDCommodityCodes
	if len(commodities) == 0 {
		commodities = defaultCommodityCodes
	}
	measures := measureFilter(entry)

	var merged []outlookRow
	var sourceURLs []string

	for _, commodity := range commodities {
		sourceURL, err := buildCommodityURL(entry, commodity)
		if err != nil {
			return nil, "", err
		}
		sourceURL += "?format=csvfilewithlabels"

		raw, err := c.download(ctx, sourceURL)
		if err != nil {
			return nil, "", fmt.Errorf("fetch oecd outlook %s: %w", commodity, err)
		}

		rows, err := parseOutlookCSV(raw, measures)
		if err != nil {
			return nil, "", err
		}
		merged = append(merged, rows...)
		sourceURLs = append(sourceURLs, sourceURL)
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("oecd ag outlook returned no rows for %s", entry.DatasetID)
	}

	sort.Slice(merged, func(i, j int) bool {
		if merged[i].CommodityCode != merged[j].CommodityCode {
			return merged[i].CommodityCode < merged[j].CommodityCode
		}
		if merged[i].MeasureCode != merged[j].MeasureCode {
			return merged[i].MeasureCode < merged[j].MeasureCode
		}
		return merged[i].Year < merged[j].Year
	})

	payload, err := json.Marshal(merged)
	if err != nil {
		return nil, "", err
	}
	return payload, strings.Join(sourceURLs, "; "), nil
}

func measureFilter(entry catalog.RegistryEntry) map[string]struct{} {
	codes := entry.OECDMeasureCodes
	if len(codes) == 0 {
		codes = defaultMeasureCodes
	}
	out := make(map[string]struct{}, len(codes))
	for _, code := range codes {
		if trimmed := strings.TrimSpace(code); trimmed != "" {
			out[trimmed] = struct{}{}
		}
	}
	return out
}

func parseOutlookCSV(raw []byte, measures map[string]struct{}) ([]outlookRow, error) {
	reader := csv.NewReader(bytes.NewReader(raw))
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read oecd outlook header: %w", err)
	}
	col := indexColumns(header)

	var rows []outlookRow
	for {
		record, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, fmt.Errorf("read oecd outlook row: %w", readErr)
		}
		if len(record) == 0 || strings.EqualFold(strings.TrimSpace(record[0]), "STRUCTURE") {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(record[0]), "DATAFLOW") {
			continue
		}

		measureCode := field(record, col, "MEASURE")
		if len(measures) > 0 {
			if _, ok := measures[measureCode]; !ok {
				continue
			}
		}

		year := field(record, col, "TIME_PERIOD")
		value := field(record, col, "OBS_VALUE")
		if year == "" || value == "" {
			continue
		}

		rows = append(rows, outlookRow{
			RefArea:       field(record, col, "REF_AREA"),
			RefAreaName:   field(record, col, "Reference area"),
			CommodityCode: field(record, col, "COMMODITY"),
			CommodityName: field(record, col, "Commodity"),
			MeasureCode:   measureCode,
			MeasureName:   field(record, col, "Measure"),
			Unit:          field(record, col, "UNIT_MEASURE"),
			UnitMult:      field(record, col, "UNIT_MULT"),
			Year:          year,
			Value:         value,
			ObsStatus:     field(record, col, "OBS_STATUS"),
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

// FlattenAgOutlook converts merged OECD-FAO outlook JSON into canonical bronze columns.
func FlattenAgOutlook(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []outlookRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse oecd ag outlook json: %w", err)
	}

	headers := []string{
		"ref_area", "ref_area_name", "commodity_code", "commodity_name",
		"measure_code", "measure_name", "unit", "unit_mult",
		"year", "value", "obs_status",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.Year) == "" || strings.TrimSpace(row.Value) == "" {
			continue
		}
		out = append(out, []string{
			row.RefArea,
			row.RefAreaName,
			row.CommodityCode,
			row.CommodityName,
			row.MeasureCode,
			row.MeasureName,
			row.Unit,
			row.UnitMult,
			row.Year,
			row.Value,
			row.ObsStatus,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no oecd ag outlook rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
