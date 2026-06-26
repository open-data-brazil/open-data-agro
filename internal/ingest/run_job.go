package ingest

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/open-data-brazil/open-data-agro/internal/alerts"
	"github.com/open-data-brazil/open-data-agro/internal/anp"
	"github.com/open-data-brazil/open-data-agro/internal/antt"
	"github.com/open-data-brazil/open-data-agro/internal/antaq"
	"github.com/open-data-brazil/open-data-agro/internal/ana"
	"github.com/open-data-brazil/open-data-agro/internal/b3"
	"github.com/open-data-brazil/open-data-agro/internal/bcb"
	"github.com/open-data-brazil/open-data-agro/internal/catalog"
	"github.com/open-data-brazil/open-data-agro/internal/cepea"
	"github.com/open-data-brazil/open-data-agro/internal/conab"
	"github.com/open-data-brazil/open-data-agro/internal/db"
	"github.com/open-data-brazil/open-data-agro/internal/dnit"
	"github.com/open-data-brazil/open-data-agro/internal/ibge"
	"github.com/open-data-brazil/open-data-agro/internal/inmet"
	"github.com/open-data-brazil/open-data-agro/internal/ipea"
	"github.com/open-data-brazil/open-data-agro/internal/mapa"
	"github.com/open-data-brazil/open-data-agro/internal/mdic"
	"github.com/open-data-brazil/open-data-agro/internal/usda"
	"github.com/open-data-brazil/open-data-agro/internal/fao"
	"github.com/open-data-brazil/open-data-agro/internal/worldbank"
	"github.com/open-data-brazil/open-data-agro/internal/noaa"
	"github.com/open-data-brazil/open-data-agro/internal/eia"
	"github.com/open-data-brazil/open-data-agro/internal/igc"
	"github.com/open-data-brazil/open-data-agro/internal/eurostat"
	"github.com/open-data-brazil/open-data-agro/internal/argentina"
	"github.com/open-data-brazil/open-data-agro/internal/un"
	"github.com/open-data-brazil/open-data-agro/internal/storage"
)

// RunOptions controls a single dataset ingest execution.
type RunOptions struct {
	DatasetID string
	DryRun    bool
	Crop      string
	FromYear  int
	ToYear    int
	Year      int
	FromDate  string
	UFs       []string
}

// RunResult summarizes an ingest attempt.
type RunResult struct {
	JobID    string
	Status   db.JobStatus
	SHA256   string
	R2Key    string
	RowCount int
	Skipped  bool
	Backend  string
}

// Runner orchestrates download, convert, upload, and audit.
type Runner struct {
	registry *catalog.Registry
	repo     *db.Repository
	store    storage.BronzeStore
	conab    *conab.Client
	anp      *anp.Client
	antt     *antt.Client
	ibge     *ibge.Client
	inmet    *inmet.Client
	bcb      *bcb.Client
	cepea    *cepea.Client
	mdic     *mdic.Client
	mapa     *mapa.Client
	b3       *b3.Client
	usda     *usda.Client
	fao      *fao.Client
	worldbank *worldbank.Client
	noaa     *noaa.Client
	eia      *eia.Client
	igc      *igc.Client
	ana      *ana.Client
	antaq    *antaq.Client
	dnit     *dnit.Client
	ipea     *ipea.Client
	eurostat *eurostat.Client
	argentina *argentina.Client
	un       *un.Client
	alerts   *alerts.Notifier
}

// NewRunner wires ingest dependencies.
func NewRunner(registry *catalog.Registry, repo *db.Repository, store storage.BronzeStore, notifier *alerts.Notifier) *Runner {
	return &Runner{
		registry: registry,
		repo:     repo,
		store:    store,
		conab:    conab.NewClient(),
		anp:      anp.NewClient(),
		antt:     antt.NewClient(),
		ibge:     ibge.NewClient(),
		inmet:    inmet.NewClient(),
		bcb:      bcb.NewClient(),
		cepea:    cepea.NewClient(),
		mdic:     mdic.NewClient(),
		mapa:     mapa.NewClient(),
		b3:       b3.NewClient(),
		usda:     usda.NewClient(),
		fao:      fao.NewClient(),
		worldbank: worldbank.NewClient(),
		noaa:     noaa.NewClient(),
		eia:      eia.NewClient(),
		igc:      igc.NewClient(),
		ana:      ana.NewClient(),
		antaq:    antaq.NewClient(),
		dnit:     dnit.NewClient(),
		ipea:     ipea.NewClient(),
		eurostat: eurostat.NewClient(),
		argentina: argentina.NewClient(),
		un:       un.NewClient(),
		alerts:   notifier,
	}
}

// Run executes one ingest job for a dataset.
func (r *Runner) Run(ctx context.Context, opts RunOptions) (*RunResult, error) {
	entry, err := r.registry.Require(opts.DatasetID)
	if err != nil {
		return nil, err
	}

	if err := r.repo.UpsertDatasetRegistry(ctx, entry); err != nil {
		return nil, fmt.Errorf("sync dataset registry: %w", err)
	}

	jobID, err := r.repo.CreateJob(ctx, opts.DatasetID, opts.DryRun)
	if err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}

	finish := func(status db.JobStatus, message *string, result *RunResult, runErr error) (*RunResult, error) {
		if finishErr := r.repo.FinishJob(ctx, jobID, status, message); finishErr != nil {
			return result, fmt.Errorf("finish job: %w (original: %v)", finishErr, runErr)
		}
		if status == db.JobFailed {
			r.alerts.Notify(ctx, slog.LevelError, fmt.Sprintf("ingest failed for %s", opts.DatasetID), "job_id", jobID, "error", derefString(message))
		}
		return result, runErr
	}

	sourceURL, err := ResolveSourceURL(entry)
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	download, err := DownloadSource(ctx, entry, r.conab, r.anp, r.antt, r.ibge, r.inmet, r.bcb, r.cepea, r.mdic, r.mapa, r.b3, r.usda, r.fao, r.worldbank, r.noaa, r.eia, r.igc, r.ana, r.antaq, r.dnit, r.ipea, r.eurostat, r.argentina, r.un, SourceOptions{
		Crop:     opts.Crop,
		FromYear: opts.FromYear,
		ToYear:   opts.ToYear,
		Year:     opts.Year,
		FromDate: opts.FromDate,
		UFs:      opts.UFs,
	})
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	ingestDate := time.Now().UTC()
	fp := NewFingerprint(download.Body, download.ContentType, download.LastModified, 0, ingestDate)

	existing, err := r.repo.FindFileBySHA256(ctx, opts.DatasetID, fp.SHA256)
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}
	if existing != nil {
		result := &RunResult{
			JobID:   jobID,
			Status:  db.JobSkipped,
			SHA256:  fp.SHA256,
			Skipped: true,
			Backend: r.store.Backend(),
		}
		if existing.R2Key != nil {
			result.R2Key = *existing.R2Key
		}
		if existing.RowCount != nil {
			result.RowCount = *existing.RowCount
		}
		_, _ = finish(db.JobSkipped, nil, result, nil)
		return result, nil
	}

	if opts.DryRun {
		result := &RunResult{
			JobID:   jobID,
			Status:  db.JobSuccess,
			SHA256:  fp.SHA256,
			Backend: r.store.Backend(),
		}
		_, _ = finish(db.JobSuccess, nil, result, nil)
		return result, nil
	}

	staged, err := StageRaw(jobID, download.Body)
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}
	defer staged.Cleanup()

	parquetBytes, rowCount, err := ConvertToParquetFromFile(entry, staged.Path)
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}
	fp.RowCount = rowCount

	bronzeKey, err := BronzeKey(opts.DatasetID, ingestDate, fp.PartID)
	if err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	if err := r.store.Put(ctx, bronzeKey, parquetBytes, "application/vnd.apache.parquet"); err != nil {
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	metadataKey, err := PartitionMetadataKey(opts.DatasetID, ingestDate)
	if err != nil {
		_ = r.store.Delete(ctx, bronzeKey)
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}
	metaSourceURL := sourceURL
	if download.SourceURL != "" {
		metaSourceURL = download.SourceURL
	}
	meta := NewBronzeMetadata(entry, fp, metaSourceURL, r.store.Backend())
	metaBytes, err := MarshalBronzeMetadata(meta)
	if err != nil {
		_ = r.store.Delete(ctx, bronzeKey)
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}
	if err := r.store.Put(ctx, metadataKey, metaBytes, "application/json"); err != nil {
		_ = r.store.Delete(ctx, bronzeKey)
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	rowCountCopy := rowCount
	fileSize := fp.FileSizeBytes
	fileRec := db.FileRecord{
		JobID:         jobID,
		DatasetID:     opts.DatasetID,
		SHA256:        fp.SHA256,
		RowCount:      &rowCountCopy,
		R2Key:         &bronzeKey,
		ContentType:   &download.ContentType,
		LastModified:  &download.LastModified,
		FileSizeBytes: &fileSize,
	}
	if err := r.repo.CreateIngestFile(ctx, fileRec); err != nil {
		if delErr := r.store.Delete(ctx, bronzeKey); delErr != nil {
			msg := fmt.Sprintf("metadata insert failed and orphan cleanup failed: %v / %v", err, delErr)
			_, _ = finish(db.JobFailed, &msg, nil, err)
			return nil, err
		}
		_ = r.store.Delete(ctx, metadataKey)
		msg := err.Error()
		_, _ = finish(db.JobFailed, &msg, nil, err)
		return nil, err
	}

	result := &RunResult{
		JobID:    jobID,
		Status:   db.JobSuccess,
		SHA256:   fp.SHA256,
		R2Key:    bronzeKey,
		RowCount: rowCount,
		Backend:  r.store.Backend(),
	}
	_, _ = finish(db.JobSuccess, nil, result, nil)
	return result, nil
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
