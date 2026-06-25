-- Placeholder schema for Phase 0.
-- Full operational migrations land in Phase 7 (07-db-postgresql).

CREATE SCHEMA IF NOT EXISTS catalog;

COMMENT ON SCHEMA catalog IS 'Operational metadata — job runs, lineage, dataset registry';
