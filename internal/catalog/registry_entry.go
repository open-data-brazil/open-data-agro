package catalog

import (
	"strings"
	"time"
)

// DatasetFormat is the source file format for a registry entry.
type DatasetFormat string

const (
	FormatCSV  DatasetFormat = "csv"
	FormatJSON DatasetFormat = "json"
	FormatXLS  DatasetFormat = "xls"
	FormatXLSX DatasetFormat = "xlsx"
	FormatTXT  DatasetFormat = "txt"
)

// RegistryEntry describes a dataset in the operational catalog.
type RegistryEntry struct {
	DatasetID       DatasetID     `json:"datasetId" yaml:"dataset_id"`
	SourceURL       string        `json:"sourceUrl" yaml:"source_url"`
	SourcePortalURL string        `json:"sourcePortalUrl,omitempty" yaml:"source_portal_url,omitempty"`
	Format          DatasetFormat `json:"format" yaml:"format"`
	Schedule     string        `json:"schedule" yaml:"schedule"`
	ConabSection string        `json:"conabSection" yaml:"conab_section"`
	PortalLabel  string        `json:"portalLabel" yaml:"portal_label"`
	Delimiter    string        `json:"delimiter,omitempty" yaml:"delimiter,omitempty"`
	XLSXSheet    string        `json:"xlsxSheet,omitempty" yaml:"xlsx_sheet,omitempty"`
	XLSXHeaderRow int          `json:"xlsxHeaderRow,omitempty" yaml:"xlsx_header_row,omitempty"`
	ANPLPCFile   string        `json:"anpLpcFile,omitempty" yaml:"anp_lpc_file,omitempty"`
	SidraTable   string        `json:"sidraTable,omitempty" yaml:"sidra_table,omitempty"`
	SidraClassification string `json:"sidraClassification,omitempty" yaml:"sidra_classification,omitempty"`
	SidraVariables []int       `json:"sidraVariables,omitempty" yaml:"sidra_variables,omitempty"`
	SidraCrops   map[string]int `json:"sidraCrops,omitempty" yaml:"sidra_crops,omitempty"`
	PeriodStart  int           `json:"periodStart,omitempty" yaml:"period_start,omitempty"`
	PeriodEnd    int           `json:"periodEnd,omitempty" yaml:"period_end,omitempty"`
	PriorityUFs  []string      `json:"priorityUfs,omitempty" yaml:"priority_ufs,omitempty"`
	ClimateVariables []string  `json:"climateVariables,omitempty" yaml:"climate_variables,omitempty"`
	RetentionNote string       `json:"retentionNote,omitempty" yaml:"retention_note,omitempty"`
	DiscoveredAt time.Time     `json:"discoveredAt" yaml:"discovered_at"`
}

// PortalURL returns the catalog portal URL for this entry.
func (e RegistryEntry) PortalURL() string {
	if strings.TrimSpace(e.SourcePortalURL) != "" {
		return strings.TrimSpace(e.SourcePortalURL)
	}
	return SourcePortalURL(e.DatasetID)
}
