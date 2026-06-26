package usda

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

const (
	wasdeESMISIndexURL = "https://esmis.nal.usda.gov/publication/world-agricultural-supply-and-demand-estimates"
	wasdeESMISBaseURL  = "https://esmis.nal.usda.gov"
)

type wasdeRow struct {
	ReportMonth string `json:"report_month"`
	Commodity   string `json:"commodity"`
	MarketYear  string `json:"market_year"`
	Attribute   string `json:"attribute"`
	Value       string `json:"value"`
	Unit        string `json:"unit"`
}

var wasdeXMLLinkPattern = regexp.MustCompile(`href="(/sites/default/release-files/[^"]+/wasde[0-9a-z]+\.xml)"`)

// FetchWASDESnapshot downloads and merges WASDE monthly XML reports from USDA ESMIS.
func (c *Client) FetchWASDESnapshot(ctx context.Context, entry catalog.RegistryEntry, fromDate string) ([]byte, string, error) {
	indexURL := strings.TrimSpace(entry.SourceURL)
	if indexURL == "" {
		indexURL = wasdeESMISIndexURL
	}

	indexBody, err := c.Download(ctx, indexURL)
	if err != nil {
		return nil, "", fmt.Errorf("fetch wasde index: %w", err)
	}

	links := discoverWASDEXMLLinks(string(indexBody))
	if len(links) == 0 {
		return nil, "", fmt.Errorf("no wasde xml links found at %s", indexURL)
	}

	minReport := resolveWASDEMinReport(entry, fromDate)
	merged := make(map[string]wasdeRow)
	fetched := 0

	for _, link := range links {
		reportMonth, ok := wasdeReportMonthFromLink(link)
		if !ok || reportMonth < minReport {
			continue
		}

		raw, err := c.Download(ctx, link)
		if err != nil {
			continue
		}

		rows, err := parseWASDEXML(raw, reportMonth)
		if err != nil {
			continue
		}
		for _, row := range rows {
			key := strings.Join([]string{row.ReportMonth, row.Commodity, row.MarketYear, row.Attribute}, "|")
			merged[key] = row
		}
		fetched++
	}

	if len(merged) == 0 {
		return nil, "", fmt.Errorf("usda wasde returned no rows for %s", entry.DatasetID)
	}

	keys := make([]string, 0, len(merged))
	for key := range merged {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	out := make([]wasdeRow, 0, len(keys))
	for _, key := range keys {
		out = append(out, merged[key])
	}

	payload, err := json.Marshal(out)
	if err != nil {
		return nil, "", err
	}

	sourceURL := fmt.Sprintf("%s#xml_reports=%d&min_report=%s", indexURL, fetched, minReport)
	return payload, sourceURL, nil
}

func discoverWASDEXMLLinks(html string) []string {
	seen := make(map[string]struct{})
	var links []string
	for _, match := range wasdeXMLLinkPattern.FindAllStringSubmatch(html, -1) {
		if len(match) < 2 {
			continue
		}
		full := wasdeESMISBaseURL + match[1]
		if _, ok := seen[full]; ok {
			continue
		}
		seen[full] = struct{}{}
		links = append(links, full)
	}
	sort.Strings(links)
	return links
}

func wasdeReportMonthFromLink(link string) (string, bool) {
	base := strings.ToLower(link)
	idx := strings.LastIndex(base, "/wasde")
	if idx < 0 {
		return "", false
	}
	name := base[idx+len("/wasde"):]
	name = strings.TrimSuffix(name, ".xml")
	name = strings.TrimSuffix(name, "v2")
	if len(name) < 4 {
		return "", false
	}
	monthPart := name[:len(name)-2]
	yearPart := name[len(name)-2:]
	month, err := strconv.Atoi(monthPart)
	if err != nil || month < 1 || month > 12 {
		return "", false
	}
	year, err := strconv.Atoi(yearPart)
	if err != nil {
		return "", false
	}
	if year < 70 {
		year += 2000
	} else {
		year += 1900
	}
	return fmt.Sprintf("%04d-%02d", year, month), true
}

func resolveWASDEMinReport(entry catalog.RegistryEntry, fromDate string) string {
	startYear := entry.PeriodStart
	if startYear == 0 {
		startYear = 2010
	}
	minReport := fmt.Sprintf("%04d-01", startYear)

	fromDate = strings.TrimSpace(fromDate)
	if fromDate != "" {
		if t, err := time.Parse("2006-01-02", fromDate); err == nil {
			candidate := t.Format("2006-01")
			if candidate > minReport {
				minReport = candidate
			}
		}
	}
	return minReport
}

func parseWASDEXML(raw []byte, reportMonth string) ([]wasdeRow, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(raw)))

	var (
		commodity  string
		marketYear string
		attribute  string
		rows       []wasdeRow
	)

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decode wasde xml: %w", err)
		}

		switch elem := tok.(type) {
		case xml.StartElement:
			switch elem.Name.Local {
			case "m1_commodity_group":
				commodity = xmlAttr(elem, "commodity1")
			case "m1_year_group":
				marketYear = strings.TrimSpace(xmlAttr(elem, "market_year1"))
			case "s3":
				attribute = normalizeWASDEAttribute(xmlAttr(elem, "attribute1"))
			case "Cell":
				value := strings.TrimSpace(xmlAttr(elem, "cell_value1"))
				if value == "" || commodity == "" || marketYear == "" || attribute == "" {
					continue
				}
				rows = append(rows, wasdeRow{
					ReportMonth: reportMonth,
					Commodity:   commodity,
					MarketYear:  marketYear,
					Attribute:   attribute,
					Value:       value,
					Unit:        "million metric tons",
				})
			}
		}
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no wasde rows in report %s", reportMonth)
	}
	return rows, nil
}

func xmlAttr(elem xml.StartElement, name string) string {
	for _, attr := range elem.Attr {
		if attr.Name.Local == name {
			return strings.TrimSpace(attr.Value)
		}
	}
	return ""
}

func normalizeWASDEAttribute(raw string) string {
	raw = strings.ReplaceAll(raw, "&#xD;&#xA;", " ")
	raw = strings.ReplaceAll(raw, "\r\n", " ")
	raw = strings.ReplaceAll(raw, "\n", " ")
	return strings.Join(strings.Fields(raw), " ")
}

// FlattenWASDE converts merged WASDE JSON into canonical bronze columns.
func FlattenWASDE(entry catalog.RegistryEntry, raw []byte) ([]string, [][]string, error) {
	var rows []wasdeRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, nil, fmt.Errorf("parse usda wasde json: %w", err)
	}

	headers := []string{
		"report_month",
		"commodity",
		"market_year",
		"attribute",
		"value",
		"unit",
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.ReportMonth) == "" || strings.TrimSpace(row.Value) == "" {
			continue
		}
		out = append(out, []string{
			row.ReportMonth,
			row.Commodity,
			row.MarketYear,
			row.Attribute,
			row.Value,
			row.Unit,
		})
	}
	if len(out) == 0 {
		return nil, nil, fmt.Errorf("no usda wasde rows to flatten for %s", entry.DatasetID)
	}
	return headers, out, nil
}
