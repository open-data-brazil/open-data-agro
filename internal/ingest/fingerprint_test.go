package ingest

import (
	"strings"
	"testing"
	"time"
)

func TestNewFingerprintUsesUUIDv7(t *testing.T) {
	t.Parallel()

	fp := NewFingerprint([]byte("hello"), "text/plain", "", 1, mustDate(t, "2026-06-25"))
	parts := strings.Split(fp.PartID, "-")
	if len(parts) < 3 {
		t.Fatalf("invalid uuid format: %q", fp.PartID)
	}
	if !strings.HasPrefix(parts[2], "7") {
		t.Fatalf("expected uuid v7 (third group starts with 7), got %q", fp.PartID)
	}
}

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
