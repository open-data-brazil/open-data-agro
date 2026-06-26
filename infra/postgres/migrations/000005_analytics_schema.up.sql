-- Phase 29 — unified analytical schema (gold mart mirror)

CREATE SCHEMA IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.sync_runs (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  finished_at TIMESTAMPTZ,
  status TEXT NOT NULL CHECK (status IN ('running', 'success', 'failed', 'partial')),
  lake_root TEXT NOT NULL,
  tables_synced INTEGER NOT NULL DEFAULT 0,
  error_message TEXT
);

CREATE TABLE IF NOT EXISTS analytics.sync_tables (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  run_id UUID NOT NULL REFERENCES analytics.sync_runs (id) ON DELETE CASCADE,
  table_name TEXT NOT NULL,
  gold_path TEXT NOT NULL,
  row_count BIGINT NOT NULL,
  min_date TEXT,
  max_date TEXT,
  parquet_bytes BIGINT,
  synced_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS sync_tables_run_id_idx ON analytics.sync_tables (run_id);
CREATE INDEX IF NOT EXISTS sync_tables_table_name_idx ON analytics.sync_tables (table_name, synced_at DESC);

CREATE OR REPLACE VIEW analytics.v_latest_sync_tables AS
SELECT DISTINCT ON (st.table_name)
  st.table_name,
  st.row_count,
  st.min_date,
  st.max_date,
  st.gold_path,
  st.synced_at,
  sr.status AS run_status,
  sr.lake_root
FROM analytics.sync_tables st
JOIN analytics.sync_runs sr ON sr.id = st.run_id
WHERE sr.status IN ('success', 'partial')
ORDER BY st.table_name, st.synced_at DESC;
