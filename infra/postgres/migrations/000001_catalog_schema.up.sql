-- catalog schema bootstrap

CREATE SCHEMA IF NOT EXISTS catalog;

COMMENT ON SCHEMA catalog IS 'Operational metadata — job runs, lineage, dataset registry';
