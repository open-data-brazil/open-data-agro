#!/usr/bin/env python3
"""Promote bronze Parquet partitions to a local or S3-backed Delta silver table."""

from __future__ import annotations

import argparse
import json
import sys
from datetime import UTC, datetime
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq
from deltalake import DeltaTable, write_deltalake


def load_metadata_map(bronze_dir: Path) -> dict[str, dict]:
    """Map bronze parquet paths to partition _metadata.json sidecars."""
    meta: dict[str, dict] = {}
    for sidecar in bronze_dir.glob("ingest_date=*/_metadata.json"):
        try:
            payload = json.loads(sidecar.read_text(encoding="utf-8"))
        except (OSError, json.JSONDecodeError):
            continue
        partition = sidecar.parent
        for parquet_file in partition.glob("part-*.parquet"):
            meta[str(parquet_file.resolve())] = payload
    return meta


def promoted_source_files(silver_uri: str, storage_options: dict[str, str] | None) -> set[str]:
    if not DeltaTable.is_deltatable(silver_uri, storage_options=storage_options):
        return set()
    dt = DeltaTable(silver_uri, storage_options=storage_options)
    table = dt.to_pyarrow_table(columns=["_source_file"])
    return {str(v) for v in table["_source_file"].to_pylist()}


def read_bronze_tables(bronze_dir: Path, dataset_id: str, skip_sources: set[str]) -> pa.Table:
    parquet_files = sorted(bronze_dir.glob("ingest_date=*/part-*.parquet"))
    if not parquet_files:
        raise FileNotFoundError(f"no bronze parquet under {bronze_dir}")

    metadata_map = load_metadata_map(bronze_dir)
    tables: list[pa.Table] = []

    for path in parquet_files:
        resolved = str(path.resolve())
        if resolved in skip_sources:
            continue
        table = pq.read_table(path)
        sidecar = metadata_map.get(resolved, {})

        ingested_at = sidecar.get("ingested_at")
        if not ingested_at:
            ingested_at = datetime.now(UTC).replace(microsecond=0).isoformat().replace("+00:00", "Z")

        n = table.num_rows
        table = table.append_column("_dataset_id", pa.array([dataset_id] * n, type=pa.string()))
        table = table.append_column("_ingested_at", pa.array([ingested_at] * n, type=pa.string()))
        table = table.append_column("_source_file", pa.array([resolved] * n, type=pa.string()))
        tables.append(table)

    if not tables:
        return pa.table({"_dataset_id": pa.array([], type=pa.string())})

    return pa.concat_tables(tables, promote_options="default")


def storage_options_from_env(storage_mode: str) -> dict[str, str] | None:
    import os

    if storage_mode == "minio":
        return {
            "AWS_ENDPOINT_URL": os.environ.get("MINIO_ENDPOINT", "http://localhost:9000"),
            "AWS_ACCESS_KEY_ID": os.environ.get("MINIO_ACCESS_KEY", "minioadmin"),
            "AWS_SECRET_ACCESS_KEY": os.environ.get("MINIO_SECRET_KEY", "minioadmin"),
            "AWS_REGION": "auto",
            "AWS_ALLOW_HTTP": "true",
            "AWS_S3_ALLOW_UNSAFE_RENAME": "true",
        }
    if storage_mode == "r2":
        endpoint = os.environ.get("R2_ENDPOINT", "")
        return {
            "AWS_ENDPOINT_URL": endpoint,
            "AWS_ACCESS_KEY_ID": os.environ.get("R2_ACCESS_KEY_ID", ""),
            "AWS_SECRET_ACCESS_KEY": os.environ.get("R2_SECRET_ACCESS_KEY", ""),
            "AWS_REGION": "auto",
        }
    return None


def promote(
    bronze_dir: Path,
    silver_dir: Path,
    dataset_id: str,
    storage_mode: str,
    min_versions: int,
) -> int:
    storage_options = storage_options_from_env(storage_mode)
    silver_uri = str(silver_dir)
    if storage_mode in ("minio", "r2"):
        bucket = (
            __import__("os").environ.get("MINIO_BUCKET" if storage_mode == "minio" else "R2_BUCKET", "open-data-agro")
        )
        table_name = silver_dir.as_posix().split("/conab/", 1)[-1]
        silver_uri = f"s3://{bucket}/silver/conab/{table_name}"

    skip_sources = promoted_source_files(silver_uri, storage_options)
    table = read_bronze_tables(bronze_dir, dataset_id, skip_sources)
    row_count = table.num_rows
    if row_count == 0:
        print(json.dumps({"row_count": 0, "silver_dir": str(silver_dir), "skipped": True}))
        return 0

    delta_log = silver_dir / "_delta_log"
    if storage_mode == "local":
        exists = delta_log.exists()
    else:
        exists = DeltaTable.is_deltatable(silver_uri, storage_options=storage_options)

    mode = "append" if exists else "overwrite"
    retention = retention_configuration(min_versions)
    write_deltalake(
        silver_uri,
        table,
        mode=mode,
        schema_mode="merge",
        storage_options=storage_options,
        configuration=retention,
    )

    return row_count


def retention_configuration(min_versions: int) -> dict[str, str]:
    """Map DELTA_MIN_VERSIONS policy to Delta table log retention."""
    days = max(7, min_versions)
    return {
        "delta.logRetentionDuration": f"interval {days} days",
        "delta.deletedFileRetentionDuration": "interval 7 days",
    }


def main() -> int:
    parser = argparse.ArgumentParser(description="Promote bronze Parquet to Delta silver")
    parser.add_argument("--bronze-dir", required=True, help="Bronze dataset directory")
    parser.add_argument("--silver-dir", required=True, help="Delta silver table directory")
    parser.add_argument("--dataset-id", required=True, help="Catalog dataset ID")
    parser.add_argument("--storage-mode", default="local", choices=["local", "minio", "r2"])
    parser.add_argument("--min-versions", type=int, default=30)
    args = parser.parse_args()

    try:
        rows = promote(
            Path(args.bronze_dir),
            Path(args.silver_dir),
            args.dataset_id,
            args.storage_mode,
            args.min_versions,
        )
    except Exception as exc:  # noqa: BLE001 — CLI boundary
        print(f"promote failed: {exc}", file=sys.stderr)
        return 1

    print(json.dumps({"row_count": rows, "silver_dir": args.silver_dir}))
    return 0


if __name__ == "__main__":
    sys.exit(main())
