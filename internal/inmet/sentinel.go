package inmet

import (
	"strings"
)

var missingSentinels = map[string]struct{}{
	"":       {},
	"NULL":   {},
	"//":     {},
	"9999":   {},
	"-9999":  {},
	"9999.0": {},
	"-9999.0": {},
}

// IsMissingValue reports whether an INMET raw cell is a documented missing sentinel.
func IsMissingValue(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	if _, ok := missingSentinels[strings.ToUpper(trimmed)]; ok {
		return true
	}
	return trimmed == "..."
}
