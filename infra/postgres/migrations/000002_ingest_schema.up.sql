-- ingestor operational tables

CREATE TABLE IF NOT EXISTS catalog.dataset_registry (
  dataset_id TEXT PRIMARY KEY,
  source_url TEXT NOT NULL,
  source_portal_url TEXT NOT NULL DEFAULT 'https://portaldeinformacoes.conab.gov.br/download-arquivos.html',
  format TEXT NOT NULL,
  schedule TEXT NOT NULL,
  conab_section TEXT,
  portal_label TEXT,
  discovered_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS catalog.ingest_jobs (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  dataset_id TEXT NOT NULL REFERENCES catalog.dataset_registry (dataset_id),
  started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  finished_at TIMESTAMPTZ,
  status TEXT NOT NULL CHECK (status IN ('running', 'success', 'failed', 'skipped')),
  error_message TEXT,
  dry_run BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS ingest_jobs_dataset_started_idx
  ON catalog.ingest_jobs (dataset_id, started_at DESC);

CREATE TABLE IF NOT EXISTS catalog.ingest_files (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  job_id UUID NOT NULL REFERENCES catalog.ingest_jobs (id),
  dataset_id TEXT NOT NULL,
  sha256 TEXT NOT NULL,
  row_count INTEGER,
  r2_key TEXT,
  content_type TEXT,
  last_modified TEXT,
  file_size_bytes BIGINT,
  discovered_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (dataset_id, sha256)
);

CREATE INDEX IF NOT EXISTS ingest_files_dataset_created_idx
  ON catalog.ingest_files (dataset_id, created_at DESC);
