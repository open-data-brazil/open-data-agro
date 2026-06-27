package bcb

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

var sgsCodePattern = regexp.MustCompile(`bcdata\.sgs\.(\d+)/dados`)

// BuildProbeURL returns an SGS JSON URL with a recent date window (required by the API).
func BuildProbeURL(entry catalog.RegistryEntry) (string, error) {
	code := entry.SGSCode
	if code == 0 {
		parsed, err := parseSGSCodeFromURL(entry.SourceURL)
		if err != nil {
			return "", err
		}
		code = parsed
	}

	end := time.Now().UTC()
	start := end.AddDate(0, -3, 0)
	return buildSGSURL(code, start.Format("02/01/2006"), end.Format("02/01/2006")), nil
}

func parseSGSCodeFromURL(raw string) (int, error) {
	match := sgsCodePattern.FindStringSubmatch(strings.TrimSpace(raw))
	if len(match) < 2 {
		return 0, fmt.Errorf("dataset missing sgs_code and unparseable source_url")
	}
	code, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("parse sgs code: %w", err)
	}
	return code, nil
}
