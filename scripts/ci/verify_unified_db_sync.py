#!/usr/bin/env python3
"""Verify PostgreSQL analytics tables match gold mart row counts."""

from __future__ import annotations

import argparse
import os
import subprocess
import sys
from pathlib import Path

import pyarrow.parquet as pq


def mart_table_name(dir_name: str) -> str:
    if not dir_name.startswith("mart_"):
        raise ValueError(dir_name)
    return dir_name.removeprefix("mart_").replace("__", "_")


def discover_marts(lake_root: Path) -> list[tuple[str, Path]]:
    gold = lake_root / "gold"
    if not gold.is_dir():
        raise FileNotFoundError(f"missing gold dir: {gold}")
    out: list[tuple[str, Path]] = []
    for entry in sorted(gold.iterdir()):
        if not entry.is_dir() or not entry.name.startswith("mart_"):
            continue
        parquet = entry / "mart.parquet"
        if parquet.is_file():
            out.append((mart_table_name(entry.name), parquet))
    return out


def pg_count(database_url: str, table_name: str) -> int:
    sql = f"SELECT COUNT(*) FROM analytics.{table_name}"
    commands = [
        ["psql", database_url, "-tAc", sql],
        [
            "docker",
            "compose",
            "exec",
            "-T",
            "postgres",
            "psql",
            "-U",
            "open_data_agro",
            "-d",
            "open_data_agro",
            "-tAc",
            sql,
        ],
    ]
    last_error = ""
    for cmd in commands:
        proc = subprocess.run(cmd, check=False, capture_output=True, text=True)
        if proc.returncode == 0:
            return int(proc.stdout.strip())
        last_error = proc.stderr.strip() or proc.stdout.strip()
    raise RuntimeError(last_error)


def main() -> int:
    parser = argparse.ArgumentParser(description="Verify analytics.* row counts vs gold marts")
    parser.add_argument("--lake-root", default=os.environ.get("LAKE_LOCAL_ROOT", "./lake"))
    parser.add_argument("--database-url", default=os.environ.get("DATABASE_URL", ""))
    args = parser.parse_args()

    if not args.database_url.strip():
        print("DATABASE_URL is required", file=sys.stderr)
        return 2

    lake_root = Path(args.lake_root).resolve()
    marts = discover_marts(lake_root)
    if not marts:
        print(f"no gold marts under {lake_root / 'gold'}", file=sys.stderr)
        return 2

    errors: list[str] = []
    for table_name, parquet_path in marts:
        gold_count = pq.read_table(parquet_path).num_rows
        try:
            pg_rows = pg_count(args.database_url, table_name)
        except RuntimeError as exc:
            errors.append(f"analytics.{table_name}: {exc}")
            continue
        if gold_count != pg_rows:
            errors.append(
                f"analytics.{table_name}: gold={gold_count} postgres={pg_rows} ({parquet_path})"
            )
        else:
            print(f"ok analytics.{table_name} rows={pg_rows}")

    if errors:
        for err in errors:
            print(err, file=sys.stderr)
        return 1

    print(f"verify_unified_db_sync: PASS ({len(marts)} marts)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
