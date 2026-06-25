package catalog

import "time"

// DatasetFormat is the source file format for a registry entry.
type DatasetFormat string

const (
	FormatCSV  DatasetFormat = "csv"
	FormatXLSX DatasetFormat = "xlsx"
	FormatTXT  DatasetFormat = "txt"
)

// RegistryEntry describes a dataset in the operational catalog.
type RegistryEntry struct {
	DatasetID    DatasetID     `json:"datasetId" yaml:"dataset_id"`
	SourceURL    string        `json:"sourceUrl" yaml:"source_url"`
	Format       DatasetFormat `json:"format" yaml:"format"`
	Schedule     string        `json:"schedule" yaml:"schedule"`
	ConabSection string        `json:"conabSection" yaml:"conab_section"`
	PortalLabel  string        `json:"portalLabel" yaml:"portal_label"`
	Delimiter    string        `json:"delimiter,omitempty" yaml:"delimiter,omitempty"`
	DiscoveredAt time.Time     `json:"discoveredAt" yaml:"discovered_at"`
}
