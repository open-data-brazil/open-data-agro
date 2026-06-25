package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parseFromFlag accepts a PAM year (2010) or ISO date (2010-01-01).
func parseFromFlag(raw string) (year int, date string, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, "", nil
	}
	if strings.Contains(raw, "-") {
		if _, parseErr := time.Parse("2006-01-02", raw); parseErr != nil {
			return 0, "", fmt.Errorf("invalid --from date %q: use YYYY-MM-DD", raw)
		}
		return 0, raw, nil
	}
	year, err = strconv.Atoi(raw)
	if err != nil || year < 1900 {
		return 0, "", fmt.Errorf("invalid --from year %q", raw)
	}
	return year, "", nil
}
