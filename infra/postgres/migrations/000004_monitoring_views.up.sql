-- monitoring views for ingestor status and ops dashboards

CREATE OR REPLACE VIEW catalog.v_latest_successful_ingest AS
SELECT DISTINCT ON (dataset_id)
  dataset_id,
  id AS job_id,
  started_at,
  finished_at,
  status,
  dry_run
FROM catalog.ingest_jobs
WHERE status = 'success'
ORDER BY dataset_id, started_at DESC;

CREATE OR REPLACE VIEW catalog.v_failed_jobs_last_7d AS
SELECT
  id,
  dataset_id,
  started_at,
  finished_at,
  status,
  error_message,
  dry_run
FROM catalog.ingest_jobs
WHERE status = 'failed'
  AND started_at >= NOW() - INTERVAL '7 days'
ORDER BY started_at DESC;
