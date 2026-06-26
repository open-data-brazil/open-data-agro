package dnit

import (
	"bytes"
	"fmt"
	"strings"
)

// StripMetadataRows removes SNV preamble lines before the CSV header row (starts with "BR").
func StripMetadataRows(raw []byte) []byte {
	lines := bytes.Split(raw, []byte("\n"))
	for i, line := range lines {
		trimmed := strings.TrimSpace(string(line))
		if trimmed == "" {
			continue
		}
		firstField := trimmed
		if idx := strings.Index(trimmed, ";"); idx >= 0 {
			firstField = trimmed[:idx]
		}
		if strings.EqualFold(strings.TrimSpace(firstField), "BR") {
			return bytes.Join(lines[i:], []byte("\n"))
		}
	}
	return raw
}

// PrepareCSV returns UTF-8 CSV bytes with metadata rows stripped.
func PrepareCSV(raw []byte) ([]byte, error) {
	stripped := StripMetadataRows(raw)
	if len(bytes.TrimSpace(stripped)) == 0 {
		return nil, fmt.Errorf("dnit csv has no data rows after stripping metadata")
	}
	return stripped, nil
}
