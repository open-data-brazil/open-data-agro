-- Promotion preview: bronze Parquet with silver metadata columns.
-- Delta write is performed by scripts/delta/promote.py (ADR 003).
-- Variables: ${bronze_uri}, ${dataset_id}
SELECT
  src.* EXCLUDE (filename),
  '${dataset_id}' AS _dataset_id,
  current_timestamp::VARCHAR AS _ingested_at,
  src.filename AS _source_file
FROM read_parquet('${bronze_uri}', filename = true) AS src;
