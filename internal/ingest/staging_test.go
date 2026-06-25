package ingest_test

import (
	"os"
	"testing"

	"github.com/open-data-brazil/open-data-agro/internal/ingest"
)

func TestStageRawWritesAndCleansUp(t *testing.T) {
	t.Parallel()

	staged, err := ingest.StageRaw("job-test", []byte("hello"))
	if err != nil {
		t.Fatalf("StageRaw: %v", err)
	}
	defer staged.Cleanup()

	data, err := os.ReadFile(staged.Path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected staged payload: %q", data)
	}

	staged.Cleanup()
	if _, err := os.Stat(staged.Path); !os.IsNotExist(err) {
		t.Fatalf("expected staged file removed, stat err=%v", err)
	}
}
