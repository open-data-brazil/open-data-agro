-- Generic promotion preview SQL (parameterized bronze → silver shape).
-- Variables: ${bronze_uri}, ${dataset_id}, ${silver_dir}
SELECT count(*) AS row_count
FROM (
  SELECT
    src.* EXCLUDE (filename),
    '${dataset_id}' AS _dataset_id,
    current_timestamp::VARCHAR AS _ingested_at,
    src.filename AS _source_file
  FROM read_parquet('${bronze_uri}', filename = true) AS src
) staged;
