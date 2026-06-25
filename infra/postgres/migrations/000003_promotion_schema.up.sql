-- processor promotion audit

CREATE TABLE IF NOT EXISTS catalog.promotion_jobs (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  dataset_id TEXT NOT NULL,
  started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  finished_at TIMESTAMPTZ,
  status TEXT NOT NULL CHECK (status IN ('running', 'success', 'failed', 'skipped', 'quality_failed')),
  row_count INTEGER,
  silver_path TEXT,
  storage_mode TEXT,
  error_message TEXT
);

CREATE INDEX IF NOT EXISTS promotion_jobs_dataset_started_idx
  ON catalog.promotion_jobs (dataset_id, started_at DESC);
