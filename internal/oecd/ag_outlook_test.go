package oecd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
)

func TestParseOutlookCSVGoldenVector(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile(filepath.Join("testdata", "ag_outlook.sample.csv"))
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}

	rows, err := parseOutlookCSV(raw, measureFilter(catalog.RegistryEntry{}))
	if err != nil {
		t.Fatalf("parseOutlookCSV: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("expected rows")
	}
	if rows[0].CommodityCode != "CPC_0141" {
		t.Fatalf("commodity: got %q", rows[0].CommodityCode)
	}
}

func TestFlattenAgOutlookGoldenVector(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID: catalog.MustParseDatasetID("oecd-fao.ag-outlook"),
	}

	raw, err := os.ReadFile(filepath.Join("testdata", "ag_outlook.sample.csv"))
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}
	parsed, err := parseOutlookCSV(raw, measureFilter(entry))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	payload, err := json.Marshal(parsed)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	headers, out, err := FlattenAgOutlook(entry, payload)
	if err != nil {
		t.Fatalf("FlattenAgOutlook: %v", err)
	}
	if len(headers) != 11 || len(out) == 0 {
		t.Fatalf("headers=%d rows=%d", len(headers), len(out))
	}
}

func TestResolveURL(t *testing.T) {
	t.Parallel()

	entry := catalog.RegistryEntry{
		DatasetID:          catalog.MustParseDatasetID("oecd-fao.ag-outlook"),
		OECDAgency:         "OECD.TAD.ATM",
		OECDDatasetID:      "DSD_AGR@DF_OUTLOOK_2024_2033",
		OECDDatasetVersion: "1.1",
		OECDDataSelection:  "BRA",
		OECDCommodityCodes: []string{"CPC_0141"},
	}

	url, err := ResolveURL(entry)
	if err != nil {
		t.Fatalf("ResolveURL: %v", err)
	}
	if url == "" || !strings.Contains(url, "sdmx.oecd.org") || !strings.Contains(url, "csvfilewithlabels") {
		t.Fatalf("unexpected url: %s", url)
	}
}
