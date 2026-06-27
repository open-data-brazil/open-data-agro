package un

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// BuildProbeURL returns a minimal Comtrade API URL (preview tier when no subscription key).
func BuildProbeURL(entry catalog.RegistryEntry) (string, error) {
	reporter := strings.TrimSpace(entry.ComtradeReporterCode)
	if reporter == "" {
		reporter = "76"
	}

	cmdCodes := entry.ComtradeCmdCodes
	if len(cmdCodes) == 0 {
		cmdCodes = []string{"1201"}
	}

	flow := "X"
	if codes := entry.ComtradeFlowCodes; len(codes) > 0 {
		flow = strings.TrimSpace(codes[0])
	}

	year := time.Now().UTC().Year() - 1
	if entry.PeriodEnd > 0 {
		year = entry.PeriodEnd
	}

	apiKey := strings.TrimSpace(os.Getenv("COMTRADE_SUBSCRIPTION_KEY"))
	return buildComtradeURL(reporter, strconv.Itoa(year), flow, cmdCodes[:1], apiKey)
}
