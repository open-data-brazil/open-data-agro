package ingest

import (
	"errors"
	"fmt"
)

// ErrProbeSkipped indicates the source health probe should not count as a link failure
// (for example when a required API key secret is not configured in CI).
var ErrProbeSkipped = errors.New("probe skipped")

// SkipProbe returns an error that ProbeCatalogEntry maps to sourceprobe.ProbeSkipped.
func SkipProbe(reason string) error {
	return fmt.Errorf("%w: %s", ErrProbeSkipped, reason)
}
