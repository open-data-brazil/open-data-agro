-- Placeholder schema for Phase 0.
-- Full operational migrations land in Phase 7 (07-db-postgresql).

CREATE SCHEMA IF NOT EXISTS catalog;

COMMENT ON SCHEMA catalog IS 'Operational metadata — job runs, lineage, dataset registry';

-- UUID policy: PostgreSQL 18.4 native uuidv7() on all PK defaults (see infra/postgres/README.md).
