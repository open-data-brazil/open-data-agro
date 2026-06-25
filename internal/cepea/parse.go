package cepea

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

var (
	dateBRPattern = regexp.MustCompile(`^\d{2}/\d{2}/\d{4}$`)
	pctPattern    = regexp.MustCompile(`[%+\s]`)
)

// Observation is one daily CEPEA indicator quote.
type Observation struct {
	Data           string `json:"data"`
	PrecoRsSc      string `json:"preco_rs_sc"`
	VariacaoDiaPct string `json:"variacao_dia_pct,omitempty"`
	PrecoUsdSc     string `json:"preco_usd_sc,omitempty"`
}

// ParseIndicatorHTML extracts daily observations for a trading location from HTML.
func ParseIndicatorHTML(htmlBody []byte, praca string) ([]Observation, error) {
	doc, err := html.Parse(bytes.NewReader(htmlBody))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	pracaKey := normalizePracaKey(praca)
	var observations []Observation
	seen := make(map[string]struct{})

	for _, table := range findTables(doc) {
		section := sectionTitleForTable(table)
		if section != "" && !pracaMatches(section, pracaKey) {
			continue
		}
		rows, err := parseTable(table)
		if err != nil {
			continue
		}
		for _, row := range rows {
			if _, ok := seen[row.Data]; ok {
				continue
			}
			if section != "" && !pracaMatches(section, pracaKey) {
				continue
			}
			seen[row.Data] = struct{}{}
			observations = append(observations, row)
		}
	}

	if len(observations) == 0 {
		return nil, fmt.Errorf("no observations found for praca %q", praca)
	}
	return observations, nil
}

func findTables(n *html.Node) []*html.Node {
	var tables []*html.Node
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "table" {
			tables = append(tables, node)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return tables
}

func sectionTitleForTable(table *html.Node) string {
	for p := table.Parent; p != nil; p = p.Parent {
		for s := p.FirstChild; s != nil; s = s.NextSibling {
			if s.Type != html.ElementNode {
				continue
			}
			tag := strings.ToLower(s.Data)
			if tag == "h1" || tag == "h2" || tag == "h3" {
				text := strings.TrimSpace(collectText(s))
				if text != "" {
					return text
				}
			}
			if tag == "div" && hasClass(s, "imagenet-table-titulo") {
				return strings.TrimSpace(collectText(s))
			}
		}
	}
	return ""
}

func parseTable(table *html.Node) ([]Observation, error) {
	headers, rows := extractRows(table)
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty table")
	}

	col := mapColumnIndexes(headers)
	if col.date < 0 || col.precoRS < 0 {
		return nil, fmt.Errorf("missing required columns in %v", headers)
	}

	out := make([]Observation, 0, len(rows))
	for _, cells := range rows {
		if col.date >= len(cells) || col.precoRS >= len(cells) {
			continue
		}
		isoDate, err := normalizeDateBR(strings.TrimSpace(cells[col.date]))
		if err != nil {
			continue
		}
		preco, err := normalizeDecimal(strings.TrimSpace(cells[col.precoRS]))
		if err != nil {
			continue
		}
		obs := Observation{
			Data:      isoDate,
			PrecoRsSc: preco,
		}
		if col.varDia >= 0 && col.varDia < len(cells) {
			if v, err := normalizePct(strings.TrimSpace(cells[col.varDia])); err == nil {
				obs.VariacaoDiaPct = v
			}
		}
		if col.precoUSD >= 0 && col.precoUSD < len(cells) {
			if v, err := normalizeDecimal(strings.TrimSpace(cells[col.precoUSD])); err == nil {
				obs.PrecoUsdSc = v
			}
		}
		out = append(out, obs)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no valid rows")
	}
	return out, nil
}

type columnIndexes struct {
	date     int
	precoRS  int
	varDia   int
	precoUSD int
}

func mapColumnIndexes(headers []string) columnIndexes {
	idx := columnIndexes{date: -1, precoRS: -1, varDia: -1, precoUSD: -1}
	for i, header := range headers {
		h := strings.ToLower(strings.TrimSpace(header))
		switch {
		case strings.Contains(h, "data") || h == "dia":
			idx.date = i
		case strings.Contains(h, "valor r") || strings.Contains(h, "valor (r$") || (strings.Contains(h, "valor") && strings.Contains(h, "r$")):
			idx.precoRS = i
		case strings.Contains(h, "var") && (strings.Contains(h, "dia") || strings.Contains(h, "/dia") || strings.Contains(h, "(%)")):
			idx.varDia = i
		case strings.Contains(h, "valor us") || strings.Contains(h, "us$"):
			idx.precoUSD = i
		}
	}
	if idx.date < 0 && len(headers) > 0 {
		idx.date = 0
	}
	if idx.precoRS < 0 && len(headers) > 1 {
		idx.precoRS = 1
	}
	if idx.varDia < 0 && len(headers) > 2 {
		idx.varDia = 2
	}
	return idx
}

func extractRows(table *html.Node) ([]string, [][]string) {
	var headers []string
	var rows [][]string
	var walk func(*html.Node, bool)
	walk = func(node *html.Node, inTable bool) {
		if node.Type == html.ElementNode && node.Data == "tr" {
			cells := rowCells(node)
			if len(cells) == 0 {
				return
			}
			if isHeaderRow(node) && len(headers) == 0 {
				headers = cells
				return
			}
			rows = append(rows, cells)
			return
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c, inTable || (node.Type == html.ElementNode && node.Data == "table"))
		}
	}
	walk(table, false)
	return headers, rows
}

func rowCells(tr *html.Node) []string {
	var cells []string
	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
			cells = append(cells, strings.TrimSpace(collectText(c)))
		}
	}
	return cells
}

func isHeaderRow(tr *html.Node) bool {
	for c := tr.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "th" {
			return true
		}
	}
	return false
}

func collectText(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			b.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.Join(strings.Fields(b.String()), " ")
}

func hasClass(n *html.Node, class string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" && strings.Contains(attr.Val, class) {
			return true
		}
	}
	return false
}

func normalizePracaKey(praca string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(praca) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func pracaMatches(section, pracaKey string) bool {
	sectionKey := normalizePracaKey(section)
	if sectionKey == "" || pracaKey == "" {
		return true
	}
	return strings.Contains(sectionKey, pracaKey) || strings.Contains(pracaKey, sectionKey)
}

func normalizeDateBR(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if !dateBRPattern.MatchString(raw) {
		return "", fmt.Errorf("invalid date %q", raw)
	}
	parts := strings.Split(raw, "/")
	return parts[2] + "-" + parts[1] + "-" + parts[0], nil
}

func normalizeDecimal(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty decimal")
	}
	raw = strings.ReplaceAll(raw, ".", "")
	raw = strings.ReplaceAll(raw, ",", ".")
	return raw, nil
}

func normalizePct(raw string) (string, error) {
	raw = pctPattern.ReplaceAllString(raw, "")
	return normalizeDecimal(raw)
}
