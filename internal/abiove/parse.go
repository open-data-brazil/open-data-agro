package abiove

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/xuri/excelize/v2"
)

var monthLabels = map[string]bool{
	"JAN": true, "FEV": true, "MAR": true, "ABR": true, "MAI": true, "JUN": true,
	"JUL": true, "AGO": true, "SET": true, "OUT": true, "NOV": true, "DEZ": true,
	"JAN-MAI": true, "JAN-MAY": true,
}

// ConvertWorkbook flattens an Abiove statistics workbook into bronze columns.
func ConvertWorkbook(entry catalog.RegistryEntry, book *excelize.File) ([]string, [][]string, error) {
	sheet := strings.TrimSpace(entry.XLSXSheet)
	if sheet == "" {
		sheet = book.GetSheetName(0)
	}
	if sheet == "" {
		return nil, nil, fmt.Errorf("workbook has no sheets")
	}

	table, err := book.GetRows(sheet)
	if err != nil {
		return nil, nil, fmt.Errorf("read sheet %q: %w", sheet, err)
	}
	if len(table) == 0 {
		return nil, nil, fmt.Errorf("sheet %q is empty", sheet)
	}

	updatedAt := findUpdatedAt(table)
	switch entry.DatasetID.String() {
	case "abiove.balanco-complexo-soja":
		return flattenBalancoExportacao(table, updatedAt)
	case "abiove.exportacoes-complexo-soja":
		return flattenWideTable(table, "Matéria-prima", updatedAt)
	case "abiove.capacidade-instalada-esmagamento":
		return flattenMonthlyBalanceRow(table, "1.3. Processamento", updatedAt)
	default:
		return nil, nil, fmt.Errorf("unsupported abiove dataset %s", entry.DatasetID)
	}
}

func flattenBalancoExportacao(table [][]string, updatedAt string) ([]string, [][]string, error) {
	start := findRowContains(table, "1.1. Exportações de soja em grão")
	if start < 0 {
		return nil, nil, fmt.Errorf("soy export section not found")
	}

	headerIdx := findRowContainsFrom(table, start, "Mês")
	if headerIdx < 0 {
		return nil, nil, fmt.Errorf("month header row not found")
	}
	baseCol := firstNonEmptyCol(table[headerIdx])

	headers := []string{"section", "row_label", "period", "metric", "value", "report_updated_at"}
	var rows [][]string
	section := "exportacoes_soja_grao"

	for i := headerIdx + 2; i < len(table); i++ {
		row := table[i]
		label := strings.TrimSpace(cellAt(row, baseCol))
		if label == "" {
			continue
		}
		if strings.HasPrefix(label, "1.2.") || strings.HasPrefix(label, "Fonte:") {
			break
		}
		if !isMonthLabel(label) && !strings.EqualFold(label, "Total ano") {
			continue
		}

		appendMetric := func(col int, metric string) {
			val := formatCell(cellAt(row, col))
			if val == "" || val == "-" {
				return
			}
			rows = append(rows, []string{section, label, label, metric, val, updatedAt})
		}
		appendMetric(baseCol+1, "valor_fob_usd_mil")
		appendMetric(baseCol+2, "peso_liquido_mil_t")
		appendMetric(baseCol+3, "preco_medio_usd_t")
	}

	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("no soy export rows parsed")
	}
	return headers, rows, nil
}

func flattenWideTable(table [][]string, headerMarker, updatedAt string) ([]string, [][]string, error) {
	headerIdx := findRowLabel(table, headerMarker)
	if headerIdx < 0 {
		return nil, nil, fmt.Errorf("header marker %q not found", headerMarker)
	}

	header := table[headerIdx]
	baseCol := firstNonEmptyCol(header)
	periods := make([]string, 0, len(header)-baseCol-1)
	for col := baseCol + 1; col < len(header); col++ {
		periods = append(periods, formatCell(cellAt(header, col)))
	}

	headers := []string{"section", "row_label", "period", "metric", "value", "report_updated_at"}
	var rows [][]string
	section := strings.ToLower(strings.ReplaceAll(headerMarker, " ", "_"))

	for i := headerIdx + 1; i < len(table); i++ {
		row := table[i]
		label := strings.TrimSpace(cellAt(row, baseCol))
		if label == "" {
			continue
		}
		if strings.HasPrefix(label, "Produção") || strings.HasPrefix(label, "Fonte:") {
			break
		}
		for col := baseCol + 1; col < len(row) && col-baseCol-1 < len(periods); col++ {
			period := periods[col-baseCol-1]
			if period == "" {
				continue
			}
			val := formatCell(cellAt(row, col))
			if val == "" {
				continue
			}
			rows = append(rows, []string{section, label, period, "volume_m3", val, updatedAt})
		}
	}

	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("no wide-table rows parsed")
	}
	return headers, rows, nil
}

func flattenMonthlyBalanceRow(table [][]string, rowPrefix, updatedAt string) ([]string, [][]string, error) {
	headerIdx := findRowContains(table, "Descriminação")
	if headerIdx < 0 {
		headerIdx = findRowContains(table, "Descrimina")
	}
	if headerIdx < 0 {
		return nil, nil, fmt.Errorf("balance header row not found")
	}

	header := table[headerIdx]
	baseCol := firstNonEmptyCol(header)
	periods := make([]string, 0, len(header)-baseCol-1)
	for col := baseCol + 1; col < len(header); col++ {
		period := formatCell(cellAt(header, col))
		if period == "" {
			continue
		}
		periods = append(periods, period)
	}

	dataIdx := findRowContains(table, rowPrefix)
	if dataIdx < 0 {
		return nil, nil, fmt.Errorf("row prefix %q not found", rowPrefix)
	}
	dataRow := table[dataIdx]
	label := strings.TrimSpace(firstNonEmptyCell(dataRow))

	headers := []string{"section", "row_label", "period", "metric", "value", "report_updated_at"}
	var rows [][]string
	section := "esmagamento_mil_t"

	for col := baseCol + 1; col < len(dataRow) && col-baseCol-1 < len(periods); col++ {
		val := formatCell(cellAt(dataRow, col))
		if val == "" {
			continue
		}
		rows = append(rows, []string{section, label, periods[col-baseCol-1], "volume_mil_t", val, updatedAt})
	}

	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("no crush capacity rows parsed")
	}
	return headers, rows, nil
}

func findUpdatedAt(table [][]string) string {
	for _, row := range table {
		for _, cell := range row {
			text := strings.TrimSpace(cell)
			if strings.HasPrefix(text, "Atualizado em:") {
				return strings.TrimSpace(strings.TrimPrefix(text, "Atualizado em:"))
			}
			if t, ok := parseExcelDate(text); ok {
				if strings.Contains(strings.ToLower(joinRow(row)), "disponíveis até") {
					return t.Format("2006-01-02")
				}
			}
		}
	}
	for _, row := range table {
		for _, cell := range row {
			if t, ok := parseExcelDate(cell); ok {
				return t.Format("2006-01-02")
			}
		}
	}
	return ""
}

func parseExcelDate(raw string) (time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, false
	}
	layouts := []string{"02/01/2006", "2006-01-02", "01/02/2006"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, true
		}
	}
	if f, err := strconv.ParseFloat(raw, 64); err == nil {
		// Excel serial date (1900-based) — approximate for metadata only.
		base := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		return base.Add(time.Duration(f*24) * time.Hour), true
	}
	return time.Time{}, false
}

func findRowContains(table [][]string, needle string) int {
	return findRowContainsFrom(table, 0, needle)
}

func findRowContainsFrom(table [][]string, start int, needle string) int {
	needle = strings.ToLower(strings.TrimSpace(needle))
	for i := start; i < len(table); i++ {
		for _, cell := range table[i] {
			if strings.Contains(strings.ToLower(strings.TrimSpace(cell)), needle) {
				return i
			}
		}
	}
	return -1
}

func findRowLabel(table [][]string, label string) int {
	want := strings.ToLower(strings.TrimSpace(label))
	for i, row := range table {
		for _, cell := range row {
			if strings.EqualFold(strings.TrimSpace(cell), want) {
				return i
			}
		}
	}
	return -1
}

func firstNonEmptyCol(row []string) int {
	for i, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return i
		}
	}
	return 0
}

func firstNonEmptyCell(row []string) string {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return cell
		}
	}
	return ""
}

func isMonthLabel(label string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(label))
	if monthLabels[normalized] {
		return true
	}
	if strings.HasPrefix(normalized, "JAN-") {
		return true
	}
	return false
}

func cellAt(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func formatCell(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "-" {
		return ""
	}
	return strings.ReplaceAll(raw, ",", ".")
}

func joinRow(row []string) string {
	return strings.Join(row, " ")
}

// OpenWorkbook opens Abiove XLSX bytes for parsing.
func OpenWorkbook(raw []byte) (*excelize.File, error) {
	return excelize.OpenReader(bytes.NewReader(raw))
}
