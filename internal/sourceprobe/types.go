package sourceprobe

import "time"

// ProbeStatus is the outcome of a single URL sample probe.
type ProbeStatus string

const (
	ProbeOK            ProbeStatus = "ok"
	ProbeUnavailable   ProbeStatus = "unavailable"
	ProbeResolveError  ProbeStatus = "resolve_error"
	ProbeSkipped       ProbeStatus = "skipped"
)

// EndpointProbe holds the result of probing one URL.
type EndpointProbe struct {
	Role           string      `json:"role"`
	URL            string      `json:"url"`
	Status         ProbeStatus `json:"status"`
	HTTPStatus     int         `json:"httpStatus,omitempty"`
	SampleBytes    int         `json:"sampleBytes,omitempty"`
	SampleSHA256   string      `json:"sampleSha256,omitempty"`
	LastModified   string      `json:"lastModified,omitempty"`
	ContentType    string      `json:"contentType,omitempty"`
	Error          string      `json:"error,omitempty"`
	Updated        bool        `json:"updated,omitempty"`
}

// DatasetProbeOutcome is the probe result for one catalog dataset.
type DatasetProbeOutcome struct {
	DatasetID  string          `json:"datasetId"`
	Endpoints  []EndpointProbe `json:"endpoints"`
	ProbedAt   time.Time       `json:"probedAt"`
	OverallOK  bool            `json:"overallOk"`
	HasUpdate  bool            `json:"hasUpdate"`
}

// Severity tracks consecutive-day source failures.
type Severity string

const (
	SeverityOK       Severity = "ok"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// HealthState persists per-dataset link health across daily runs.
type HealthState struct {
	DatasetID              string   `json:"datasetId"`
	Endpoints              []string `json:"endpoints"`
	ConsecutiveFailureDays int      `json:"consecutiveFailureDays"`
	FirstFailureDate       string   `json:"firstFailureDate,omitempty"`
	LastFailureDate        string   `json:"lastFailureDate,omitempty"`
	LastSuccessDate        string   `json:"lastSuccessDate,omitempty"`
	Severity               Severity `json:"severity"`
	Message                string   `json:"message"`
	LastSampleSHA256       string   `json:"lastSampleSha256,omitempty"`
	LastModified           string   `json:"lastModified,omitempty"`
}

// SourceAlert is a human-facing alert for a degraded dataset.
type SourceAlert struct {
	DatasetID              string   `json:"datasetId"`
	Severity               Severity `json:"severity"`
	Message                string   `json:"message"`
	Endpoints              []string `json:"endpoints"`
	ConsecutiveFailureDays int      `json:"consecutiveFailureDays"`
}

// ReportSummary aggregates a daily probe run.
type ReportSummary struct {
	ExecutedAt      time.Time `json:"executedAt"`
	RunDate         string    `json:"runDate"`
	TotalDatasets   int       `json:"totalDatasets"`
	OKCount         int       `json:"okCount"`
	WarningCount    int       `json:"warningCount"`
	CriticalCount   int       `json:"criticalCount"`
	UpdatedCount    int       `json:"updatedCount"`
	DeprecatedCount int       `json:"deprecatedCount"`
}

// DailyReport is written to latest.json and daily archive.
type DailyReport struct {
	Summary      ReportSummary         `json:"summary"`
	Outcomes     []DatasetProbeOutcome `json:"outcomes"`
	Alerts       []SourceAlert         `json:"alerts"`
	Updated      []string              `json:"updated"`
	Deprecated   []string              `json:"deprecated"`
	CommitMessage string               `json:"commitMessage"`
}
