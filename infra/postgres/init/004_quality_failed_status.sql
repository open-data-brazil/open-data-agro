-- Phase 6 — promotion quality gate status

ALTER TABLE catalog.promotion_jobs
  DROP CONSTRAINT IF EXISTS promotion_jobs_status_check;

ALTER TABLE catalog.promotion_jobs
  ADD CONSTRAINT promotion_jobs_status_check
  CHECK (status IN ('running', 'success', 'failed', 'skipped', 'quality_failed'));
