-- Smoke: count bronze Parquet rows for a dataset (local or s3:// URI).
-- Variables: ${bronze_uri}
SELECT count(*) AS row_count
FROM read_parquet('${bronze_uri}');
