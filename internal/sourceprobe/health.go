package sourceprobe

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func isPreviousCalendarDay(previousDate, currentDate string) bool {
	previous, err1 := time.Parse("2006-01-02", previousDate)
	current, err2 := time.Parse("2006-01-02", currentDate)
	if err1 != nil || err2 != nil {
		return false
	}
	diff := current.Sub(previous)
	return diff >= 24*time.Hour && diff < 48*time.Hour
}

func buildFailureMessage(consecutiveDays int) string {
	if consecutiveDays >= 2 {
		return "Consultation link deprecated — official source unreachable for 2 or more consecutive days."
	}
	return "Possible link deprecation — official source unreachable after retries."
}

// ApplyOutcome updates health state from a probe outcome.
func ApplyOutcome(previous *HealthState, outcome DatasetProbeOutcome, runDate string) HealthState {
	endpoints := endpointURLs(outcome)
	sourceProbe := sourceEndpoint(outcome)

	if sourceProbe != nil && sourceProbe.Status == ProbeOK {
		state := HealthState{
			DatasetID:        outcome.DatasetID,
			Endpoints:        endpoints,
			LastSuccessDate:  runDate,
			Severity:         SeverityOK,
			Message:          "Official source responded successfully (sample probe).",
			LastSampleSHA256: sourceProbe.SampleSHA256,
			LastModified:     sourceProbe.LastModified,
		}
		return state
	}

	errMsg := "source unavailable"
	if sourceProbe != nil && sourceProbe.Error != "" {
		errMsg = sourceProbe.Error
	} else if sourceProbe != nil && sourceProbe.Status == ProbeResolveError {
		errMsg = sourceProbe.Error
	}

	if previous != nil && previous.LastFailureDate == runDate {
		severity := SeverityWarning
		if previous.ConsecutiveFailureDays >= 2 {
			severity = SeverityCritical
		}
		return HealthState{
			DatasetID:              outcome.DatasetID,
			Endpoints:              endpoints,
			ConsecutiveFailureDays: previous.ConsecutiveFailureDays,
			FirstFailureDate:       previous.FirstFailureDate,
			LastFailureDate:        runDate,
			LastSuccessDate:        previous.LastSuccessDate,
			Severity:               severity,
			Message:                buildFailureMessage(previous.ConsecutiveFailureDays),
			LastSampleSHA256:       previous.LastSampleSHA256,
			LastModified:           previous.LastModified,
		}
	}

	consecutive := 1
	firstFailure := runDate
	if previous != nil && previous.LastFailureDate != "" && isPreviousCalendarDay(previous.LastFailureDate, runDate) {
		consecutive = previous.ConsecutiveFailureDays + 1
		firstFailure = previous.FirstFailureDate
		if firstFailure == "" {
			firstFailure = runDate
		}
	}

	severity := SeverityWarning
	if consecutive >= 2 {
		severity = SeverityCritical
	}

	return HealthState{
		DatasetID:              outcome.DatasetID,
		Endpoints:              endpoints,
		ConsecutiveFailureDays: consecutive,
		FirstFailureDate:       firstFailure,
		LastFailureDate:        runDate,
		LastSuccessDate:        previousLastSuccess(previous),
		Severity:               severity,
		Message:                buildFailureMessage(consecutive) + " (" + errMsg + ")",
		LastSampleSHA256:       previousSampleHash(previous),
		LastModified:           previousModified(previous),
	}
}

func previousLastSuccess(previous *HealthState) string {
	if previous == nil {
		return ""
	}
	return previous.LastSuccessDate
}

func previousSampleHash(previous *HealthState) string {
	if previous == nil {
		return ""
	}
	return previous.LastSampleSHA256
}

func previousModified(previous *HealthState) string {
	if previous == nil {
		return ""
	}
	return previous.LastModified
}

func endpointURLs(outcome DatasetProbeOutcome) []string {
	urls := make([]string, 0, len(outcome.Endpoints))
	for _, ep := range outcome.Endpoints {
		if ep.URL != "" {
			urls = append(urls, ep.URL)
		}
	}
	return urls
}

func sourceEndpoint(outcome DatasetProbeOutcome) *EndpointProbe {
	for i := range outcome.Endpoints {
		if outcome.Endpoints[i].Role == "source" {
			return &outcome.Endpoints[i]
		}
	}
	if len(outcome.Endpoints) > 0 {
		return &outcome.Endpoints[0]
	}
	return nil
}

// HealthAlert converts a non-ok health state to an alert.
func HealthAlert(state HealthState) (SourceAlert, bool) {
	if state.Severity == SeverityOK {
		return SourceAlert{}, false
	}
	return SourceAlert{
		DatasetID:              state.DatasetID,
		Severity:               state.Severity,
		Message:                state.Message,
		Endpoints:              state.Endpoints,
		ConsecutiveFailureDays: state.ConsecutiveFailureDays,
	}, true
}

// ReadHealthState loads persisted health for a dataset.
func ReadHealthState(dir, datasetID string) (*HealthState, error) {
	path := filepath.Join(dir, datasetID+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var state HealthState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parse health %s: %w", path, err)
	}
	return &state, nil
}

// WriteHealthState persists health for a dataset.
func WriteHealthState(dir string, state HealthState) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(dir, state.DatasetID+".json"), data, 0o644)
}
