package anp

import "strings"

const ethanolProductNeedle = "ETANOL"

// FilterEthanolPrecos keeps LPC municipal average rows whose PRODUTO is ethanol.
func FilterEthanolPrecos(headers []string, rows [][]string) ([]string, [][]string) {
	prodIdx := productColumnIndex(headers)
	if prodIdx < 0 {
		return headers, rows
	}

	filtered := make([][]string, 0, len(rows))
	for _, row := range rows {
		if prodIdx >= len(row) {
			continue
		}
		if strings.Contains(strings.ToUpper(strings.TrimSpace(row[prodIdx])), ethanolProductNeedle) {
			filtered = append(filtered, row)
		}
	}
	return headers, filtered
}

func productColumnIndex(headers []string) int {
	for i, header := range headers {
		normalized := strings.ToUpper(strings.TrimSpace(header))
		if normalized == "PRODUTO" {
			return i
		}
	}
	return -1
}
