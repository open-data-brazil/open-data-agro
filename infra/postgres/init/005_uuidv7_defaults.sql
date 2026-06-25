-- Upgrade path: switch UUID defaults from gen_random_uuid() (v4) to native uuidv7().
-- Requires PostgreSQL 18.4+ (Docker image postgres:18.4-bookworm).

ALTER TABLE catalog.ingest_jobs
  ALTER COLUMN id SET DEFAULT uuidv7();

ALTER TABLE catalog.ingest_files
  ALTER COLUMN id SET DEFAULT uuidv7();

ALTER TABLE catalog.promotion_jobs
  ALTER COLUMN id SET DEFAULT uuidv7();
