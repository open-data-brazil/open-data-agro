package mdic

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// BuildProbeRequest returns a minimal Comex Stat POST body and endpoint URL.
func BuildProbeRequest(entry catalog.RegistryEntry) ([]byte, string, error) {
	flow := strings.TrimSpace(entry.ComexFlow)
	if flow == "" {
		flow = "export"
	}

	ncms := entry.ComexNCMs
	if len(ncms) == 0 {
		ncms = []string{"12019000"}
	}

	details := entry.ComexDetails
	if len(details) == 0 {
		details = []string{"ncm"}
	}

	metrics := entry.ComexMetrics
	if len(metrics) == 0 {
		if strings.EqualFold(flow, "import") {
			metrics = []string{"metricCIF", "metricKG"}
		} else {
			metrics = []string{"metricFOB", "metricKG"}
		}
	}

	year := time.Now().UTC().Year() - 1
	periodFrom := fmt.Sprintf("%04d-01", year)
	periodTo := fmt.Sprintf("%04d-01", year)

	body, err := json.Marshal(comexRequest{
		Flow:        flow,
		MonthDetail: true,
		Period:      comexPeriod{From: periodFrom, To: periodTo},
		Filters:     []comexFilter{{Filter: "ncm", Values: ncms[:1]}},
		Details:     details,
		Metrics:     metrics,
	})
	if err != nil {
		return nil, "", err
	}

	return body, comexAPIBase + generalEndpoint, nil
}
