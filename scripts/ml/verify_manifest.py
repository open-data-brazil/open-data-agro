#!/usr/bin/env python3
"""Verify ML export manifest against parquet (Phase 30 B7)."""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import sys
from pathlib import Path

import pyarrow.parquet as pq

ROOT = Path(__file__).resolve().parents[2]


def _schema_hash(column_names: list[str]) -> str:
    payload = "|".join(column_names).encode("utf-8")
    return hashlib.sha256(payload).hexdigest()[:16]


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--export-dir",
        default=os.environ.get("ML_EXPORT_DIR", str(ROOT / ".local" / "ml" / "export")),
    )
    args = parser.parse_args()
    export_dir = Path(args.export_dir)
    manifest_path = export_dir / "manifest.json"
    parquet_path = export_dir / "soy_daily_training.parquet"

    errors: list[str] = []
    if not manifest_path.is_file():
        errors.append(f"missing {manifest_path}")
    if not parquet_path.is_file():
        errors.append(f"missing {parquet_path}")
    if errors:
        print("verify_manifest: FAIL", file=sys.stderr)
        for e in errors:
            print(f"  - {e}", file=sys.stderr)
        return 1

    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    table = pq.read_table(parquet_path)
    cols = table.column_names
    actual_hash = _schema_hash(cols)

    if manifest.get("row_count") != table.num_rows:
        errors.append(
            f"row_count mismatch: manifest={manifest.get('row_count')} parquet={table.num_rows}"
        )
    if manifest.get("schema_hash") != actual_hash:
        errors.append(
            f"schema_hash mismatch: manifest={manifest.get('schema_hash')} actual={actual_hash}"
        )
    if manifest.get("columns") != cols:
        errors.append("columns list mismatch vs parquet")
    if "y_ret_30d" not in cols or "split" not in cols:
        errors.append("missing required training columns y_ret_30d / split")

    if errors:
        print("verify_manifest: FAIL", file=sys.stderr)
        for e in errors:
            print(f"  - {e}", file=sys.stderr)
        return 1

    print(
        f"verify_manifest: PASS rows={table.num_rows} schema_hash={actual_hash} "
        f"date={manifest.get('date_min')}..{manifest.get('date_max')}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
