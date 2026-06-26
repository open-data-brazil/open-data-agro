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
	SGSCode      int           `json:"sgsCode,omitempty" yaml:"sgs_code,omitempty"`
	CepeaProductSlug string    `json:"cepeaProductSlug,omitempty" yaml:"cepea_product_slug,omitempty"`
	CepeaPraca   string        `json:"cepeaPraca,omitempty" yaml:"cepea_praca,omitempty"`
	StartDate    string        `json:"startDate,omitempty" yaml:"start_date,omitempty"`
	Frequency    string        `json:"frequency,omitempty" yaml:"frequency,omitempty"`
	License      string        `json:"license,omitempty" yaml:"license,omitempty"`
	FonteTipo    string        `json:"fonteTipo,omitempty" yaml:"fonte_tipo,omitempty"`
	ComexFlow    string        `json:"comexFlow,omitempty" yaml:"comex_flow,omitempty"`
	ComexNCMs    []string      `json:"comexNcms,omitempty" yaml:"comex_ncms,omitempty"`
	ComexDetails []string      `json:"comexDetails,omitempty" yaml:"comex_details,omitempty"`
	ComexMetrics []string      `json:"comexMetrics,omitempty" yaml:"comex_metrics,omitempty"`
	CKANPackageID string       `json:"ckanPackageId,omitempty" yaml:"ckan_package_id,omitempty"`
	CKANResourceFormat string  `json:"ckanResourceFormat,omitempty" yaml:"ckan_resource_format,omitempty"`
	CKANResourceNameContains string `json:"ckanResourceNameContains,omitempty" yaml:"ckan_resource_name_contains,omitempty"`
	B3FilePrefix string       `json:"b3FilePrefix,omitempty" yaml:"b3_file_prefix,omitempty"`
	B3CommodityPrefix string  `json:"b3CommodityPrefix,omitempty" yaml:"b3_commodity_prefix,omitempty"`
	PSDCommodityCode string   `json:"psdCommodityCode,omitempty" yaml:"psd_commodity_code,omitempty"`
	PSDCommoditySlug string   `json:"psdCommoditySlug,omitempty" yaml:"psd_commodity_slug,omitempty"`
	FAOBulkURL       string   `json:"faoBulkUrl,omitempty" yaml:"fao_bulk_url,omitempty"`
	FAOBulkCSV       string   `json:"faoBulkCsv,omitempty" yaml:"fao_bulk_csv,omitempty"`
	FAOItemCodes     []string `json:"faoItemCodes,omitempty" yaml:"fao_item_codes,omitempty"`
	FAOElementCodes  []string `json:"faoElementCodes,omitempty" yaml:"fao_element_codes,omitempty"`
	WorldBankPinkSheetURL   string   `json:"worldbankPinkSheetUrl,omitempty" yaml:"worldbank_pink_sheet_url,omitempty"`
	WorldBankPinkSheetSheet string   `json:"worldbankPinkSheetSheet,omitempty" yaml:"worldbank_pink_sheet_sheet,omitempty"`
	WorldBankSeriesNames    []string `json:"worldbankSeriesNames,omitempty" yaml:"worldbank_series_names,omitempty"`
	NOAAIndexURL            string   `json:"noaaIndexUrl,omitempty" yaml:"noaa_index_url,omitempty"`
	EIASeriesIDs            []string `json:"eiaSeriesIds,omitempty" yaml:"eia_series_ids,omitempty"`
	DiscoveredAt time.Time     `json:"discoveredAt" yaml:"discovered_at"`
}

// PortalURL returns the catalog portal URL for this entry.
func (e RegistryEntry) PortalURL() string {
	if strings.TrimSpace(e.SourcePortalURL) != "" {
		return strings.TrimSpace(e.SourcePortalURL)
	}
	return SourcePortalURL(e.DatasetID)
}
