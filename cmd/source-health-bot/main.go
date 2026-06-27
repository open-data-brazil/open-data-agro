package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/ingest"
	"github.com/open-data-brazil/open-data-agro/internal/sourceprobe"
)

func main() {
	reportDir := flag.String("report-dir", "data/source-health-reports", "directory for probe reports")
	concurrency := flag.Int("concurrency", 6, "parallel probe workers")
	runDate := flag.String("run-date", "", "run date YYYY-MM-DD (default: today UTC)")
	flag.Parse()

	reg, err := catalog.LoadDefaultRegistry()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	outcomes := ingest.ProbeAll(ctx, reg, *concurrency)
	report, err := sourceprobe.Run(ctx, outcomes, sourceprobe.RunOptions{
		ReportDir:   *reportDir,
		Concurrency: *concurrency,
		RunDate:     *runDate,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("source-health-bot: probed %d datasets — ok=%d warning=%d critical=%d updated=%d\n",
		report.Summary.TotalDatasets,
		report.Summary.OKCount,
		report.Summary.WarningCount,
		report.Summary.CriticalCount,
		report.Summary.UpdatedCount,
	)
	fmt.Printf("reports: %s/latest.json\n", *reportDir)
	fmt.Printf("suggested commit: %s\n", report.CommitMessage)

	if report.Summary.CriticalCount > 0 {
		os.Exit(2)
	}
}
