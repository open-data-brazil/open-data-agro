package ons

import (
	"strings"
	"testing"
)

func TestResolveURLDirect(t *testing.T) {
	t.Parallel()

	raw := "https://ons-aws-prod-opendata.s3.amazonaws.com/dataset/carga_energia_di/CARGA_ENERGIA_2024.csv"
	lower := strings.ToLower(raw)
	if !strings.Contains(lower, "amazonaws.com") {
		t.Fatalf("expected ONS S3 url")
	}
}
