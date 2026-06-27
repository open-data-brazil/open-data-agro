package sourceprobe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProbeURLReadsSampleOnly(t *testing.T) {
	var receivedRange string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedRange = r.Header.Get("Range")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Last-Modified", "Mon, 01 Jan 2024 00:00:00 GMT")
		_, _ = w.Write([]byte("hello-world-sample-data-for-probe"))
	}))
	defer server.Close()

	client := NewClient()
	result, err := client.ProbeURL(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("ProbeURL: %v", err)
	}
	if receivedRange != "" {
		t.Fatalf("expected GET without Range header, got %q", receivedRange)
	}
	if result.SampleBytes == 0 {
		t.Fatal("expected non-zero sample bytes")
	}
	if result.SampleSHA256 == "" {
		t.Fatal("expected sample hash")
	}
	if result.LastModified == "" {
		t.Fatal("expected Last-Modified header")
	}
}

func TestApplyOutcomeEscalatesOnConsecutiveDays(t *testing.T) {
	previous := &HealthState{
		DatasetID:              "conab.test",
		ConsecutiveFailureDays: 1,
		FirstFailureDate:       "2026-06-25",
		LastFailureDate:        "2026-06-25",
		Severity:               SeverityWarning,
	}
	outcome := DatasetProbeOutcome{
		DatasetID: "conab.test",
		Endpoints: []EndpointProbe{{
			Role:   "source",
			URL:    "https://example.test/missing.txt",
			Status: ProbeUnavailable,
			Error:  "unexpected status 404",
		}},
	}

	state := ApplyOutcome(previous, outcome, "2026-06-26")
	if state.Severity != SeverityCritical {
		t.Fatalf("expected critical, got %s", state.Severity)
	}
	if state.ConsecutiveFailureDays != 2 {
		t.Fatalf("expected 2 consecutive days, got %d", state.ConsecutiveFailureDays)
	}
}

func TestBuildCommitMessageAllOK(t *testing.T) {
	msg := buildCommitMessage(
		mustParseTime("2026-06-26T10:00:00Z"),
		nil,
		nil,
		ReportSummary{TotalDatasets: 97},
	)
	if msg == "" {
		t.Fatal("expected commit message")
	}
}

func mustParseTime(value string) (t time.Time) {
	t, _ = time.Parse(time.RFC3339, value)
	return t
}
