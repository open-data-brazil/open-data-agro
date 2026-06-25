package catalog

import (
	"strings"
	"time"
)

// DatasetFormat is the source file format for a registry entry.
type DatasetFormat string

const (
	FormatCSV  DatasetFormat = "csv"
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
	DiscoveredAt time.Time     `json:"discoveredAt" yaml:"discovered_at"`
}

// PortalURL returns the catalog portal URL for this entry.
func (e RegistryEntry) PortalURL() string {
	if strings.TrimSpace(e.SourcePortalURL) != "" {
		return strings.TrimSpace(e.SourcePortalURL)
	}
	return SourcePortalURL(e.DatasetID)
}
