package sourceprobe

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const dailyRetentionDays = 90

// RunOptions configures a daily source health probe run.
type RunOptions struct {
	ReportDir   string
	Concurrency int
	RunDate     string
}

// Run processes probe outcomes and writes reports.
func Run(ctx context.Context, outcomes []DatasetProbeOutcome, opts RunOptions) (*DailyReport, error) {
	if opts.ReportDir == "" {
		opts.ReportDir = "data/source-health-reports"
	}
	if opts.Concurrency <= 0 {
		opts.Concurrency = 6
	}
	if opts.RunDate == "" {
		opts.RunDate = time.Now().UTC().Format("2006-01-02")
	}

	healthDir := filepath.Join(opts.ReportDir, "source-health")
	outcomeDir := filepath.Join(opts.ReportDir, "probe-outcomes")

	_ = ctx
	states := make([]HealthState, 0, len(outcomes))
	alerts := make([]SourceAlert, 0)
	updated := make([]string, 0)
	deprecated := make([]string, 0)

	for i := range outcomes {
		previous, err := ReadHealthState(healthDir, outcomes[i].DatasetID)
		if err != nil {
			return nil, err
		}

		detectUpdate(previous, &outcomes[i])
		if outcomes[i].HasUpdate {
			updated = append(updated, outcomes[i].DatasetID)
		}

		state := ApplyOutcome(previous, outcomes[i], opts.RunDate)
		if err := WriteHealthState(healthDir, state); err != nil {
			return nil, err
		}
		states = append(states, state)

		if alert, ok := HealthAlert(state); ok {
			alerts = append(alerts, alert)
			if state.Severity == SeverityCritical {
				deprecated = append(deprecated, outcomes[i].DatasetID)
			}
		}

		if err := writeProbeOutcome(outcomeDir, outcomes[i]); err != nil {
			return nil, err
		}
	}

	report := buildReport(outcomes, alerts, updated, deprecated, opts.RunDate)
	if err := writeReports(opts.ReportDir, report); err != nil {
		return nil, err
	}
	if err := pruneDailyArchives(filepath.Join(opts.ReportDir, "daily")); err != nil {
		return nil, err
	}

	_ = states
	return report, nil
}

func detectUpdate(previous *HealthState, outcome *DatasetProbeOutcome) {
	source := sourceEndpoint(*outcome)
	if source == nil || source.Status != ProbeOK {
		return
	}
	if previous == nil || previous.LastSampleSHA256 == "" {
		return
	}
	if source.SampleSHA256 != previous.LastSampleSHA256 {
		source.Updated = true
		outcome.HasUpdate = true
	}
}

func buildReport(
	outcomes []DatasetProbeOutcome,
	alerts []SourceAlert,
	updated, deprecated []string,
	runDate string,
) *DailyReport {
	executedAt := time.Now().UTC()
	summary := ReportSummary{
		ExecutedAt:    executedAt,
		RunDate:       runDate,
		TotalDatasets: len(outcomes),
		UpdatedCount:  len(updated),
	}

	for _, alert := range alerts {
		switch alert.Severity {
		case SeverityWarning:
			summary.WarningCount++
		case SeverityCritical:
			summary.CriticalCount++
		}
	}
	summary.DeprecatedCount = len(deprecated)
	summary.OKCount = summary.TotalDatasets - summary.WarningCount - summary.CriticalCount

	commitMessage := buildCommitMessage(executedAt, deprecated, updated, summary)

	return &DailyReport{
		Summary:       summary,
		Outcomes:      outcomes,
		Alerts:        alerts,
		Updated:       updated,
		Deprecated:    deprecated,
		CommitMessage: commitMessage,
	}
}

func buildCommitMessage(executedAt time.Time, deprecated, updated []string, summary ReportSummary) string {
	stamp := executedAt.UTC().Format(time.RFC3339)
	if len(deprecated) == 0 && len(updated) == 0 {
		return fmt.Sprintf("chore(source-health): daily probe %s — all %d links OK", stamp, summary.TotalDatasets)
	}
	parts := []string{fmt.Sprintf("chore(source-health): daily probe %s", stamp)}
	if len(deprecated) > 0 {
		parts = append(parts, fmt.Sprintf("%d deprecated (2+ days unreachable)", len(deprecated)))
	}
	if summary.WarningCount > 0 && len(deprecated) == 0 {
		parts = append(parts, fmt.Sprintf("%d unreachable (day 1 warning)", summary.WarningCount))
	}
	if len(updated) > 0 {
		parts = append(parts, fmt.Sprintf("%d updated samples", len(updated)))
	}
	return strings.Join(parts, " — ")
}

func writeProbeOutcome(dir string, outcome DatasetProbeOutcome) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(outcome, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, outcome.DatasetID+".json"), data, 0o644)
}

func writeReports(reportDir string, report *DailyReport) error {
	if err := os.MkdirAll(reportDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(reportDir, "daily"), 0o755); err != nil {
		return err
	}

	latestPath := filepath.Join(reportDir, "latest.json")
	if err := writeJSON(latestPath, report); err != nil {
		return err
	}

	dailyPath := filepath.Join(reportDir, "daily", report.Summary.RunDate+".json")
	if err := writeJSON(dailyPath, report); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(reportDir, "job-summary.md"), []byte(renderJobSummary(report)), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(reportDir, "CRITICAL-ALERTS.md"), []byte(renderCriticalAlerts(report)), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(reportDir, "pr-body.md"), []byte(renderPRBody(report)), 0o644); err != nil {
		return err
	}

	docsPath := filepath.Join("docs", "SOURCE-HEALTH.md")
	return os.WriteFile(docsPath, []byte(renderSourceHealthDoc(report)), 0o644)
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func pruneDailyArchives(dailyDir string) error {
	cutoff := time.Now().UTC().AddDate(0, 0, -dailyRetentionDays)
	entries, err := os.ReadDir(dailyDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		dateStr := strings.TrimSuffix(entry.Name(), ".json")
		fileDate, parseErr := time.Parse("2006-01-02", dateStr)
		if parseErr != nil {
			continue
		}
		if fileDate.Before(cutoff) {
			_ = os.Remove(filepath.Join(dailyDir, entry.Name()))
		}
	}
	return nil
}
