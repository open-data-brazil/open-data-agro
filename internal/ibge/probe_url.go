package ibge

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

// BuildProbeURL returns a minimal SIDRA URL that mirrors ingest fetch parameters.
func BuildProbeURL(entry catalog.RegistryEntry) (string, error) {
	table := strings.TrimSpace(entry.SidraTable)
	if table == "" {
		return "", fmt.Errorf("dataset %s missing sidra_table", entry.DatasetID)
	}

	id := entry.DatasetID.String()
	if strings.HasPrefix(id, "ibge.lspa-") {
		classification := strings.TrimSpace(entry.SidraClassification)
		if classification == "" || len(entry.SidraVariables) == 0 {
			return "", fmt.Errorf("dataset %s missing lspa probe fields", entry.DatasetID)
		}
		cropCode := firstSidraCropCode(entry)
		if cropCode == "" {
			return "", fmt.Errorf("dataset %s missing sidra_crops", entry.DatasetID)
		}
		year := probeYear(entry)
		vars := formatVariables(entry.SidraVariables[:1])
		period := fmt.Sprintf("%d12", year)
		return buildLSPAURL(table, []string{"n3 11"}, period, vars, classification, cropCode), nil
	}

	if strings.HasPrefix(id, "ibge.pevs-") {
		if len(entry.SidraVariables) == 0 {
			return "", fmt.Errorf("dataset %s missing sidra_variables", entry.DatasetID)
		}
		year := probeYear(entry)
		vars := formatVariables(entry.SidraVariables[:1])
		return buildPEVSURL(table, strconv.Itoa(year), vars), nil
	}

	if strings.HasPrefix(id, "ibge.ppm-") {
		variables := entry.SidraVariables
		if len(variables) == 0 {
			return "", fmt.Errorf("dataset %s missing sidra_variables", entry.DatasetID)
		}
		year := probeYear(entry)
		varsProbe := formatVariables(variables[:1])
		if varsProbe == "" {
			return "", fmt.Errorf("dataset %s missing sidra_variables", entry.DatasetID)
		}
		ufBatch := []string{defaultUFChunks[0][0]}
		return buildPPMURL(table, ufBatch, year, varsProbe), nil
	}

	if strings.HasPrefix(id, "ibge.censo-agro-") {
		year := entry.PeriodEnd
		if year == 0 {
			year = 2017
		}
		return sidraValuesBase + fmt.Sprintf("/t/%s/n3/all/p/%d/v/all", table, year), nil
	}

	if strings.HasPrefix(id, "ibge.pnad-continua-") {
		variables := formatVariables(entry.SidraVariables)
		if variables == "" {
			variables = "all"
		}
		return buildPNADURL(table, defaultUFChunks[0], "last 1", variables), nil
	}

	if entry.SidraClassification == "" || len(entry.SidraVariables) == 0 {
		return "", fmt.Errorf("dataset %s missing sidra probe fields", entry.DatasetID)
	}

	year := probeYear(entry)

	cropCode := firstSidraCropCode(entry)
	if cropCode == "" {
		return "", fmt.Errorf("dataset %s missing sidra_crops", entry.DatasetID)
	}

	if len(entry.SidraVariables) == 0 {
		return "", fmt.Errorf("dataset %s missing sidra_variables", entry.DatasetID)
	}
	vars := formatVariables(entry.SidraVariables[:1])
	ufBatch := []string{defaultUFChunks[0][0]}
	return buildSIDRAURL(table, ufBatch, year, vars, entry.SidraClassification, cropCode), nil
}

func firstSidraCropCode(entry catalog.RegistryEntry) string {
	for _, code := range entry.SidraCrops {
		return strconv.Itoa(code)
	}
	return ""
}

func probeYear(entry catalog.RegistryEntry) int {
	year := entry.PeriodEnd
	if year == 0 {
		year = time.Now().UTC().Year() - 1
	}
	maxYear := time.Now().UTC().Year() - 1
	if year > maxYear {
		return maxYear
	}
	return year
}
