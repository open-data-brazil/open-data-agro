#!/usr/bin/env python3
"""Spot-check PostgreSQL analytics tables for row counts and date ranges."""

from __future__ import annotations

import argparse
import os
import subprocess
import sys


def run_psql(database_url: str, sql: str) -> str:
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
            return proc.stdout.strip()
        last_error = proc.stderr.strip() or proc.stdout.strip()
    raise RuntimeError(last_error)


def main() -> int:
    parser = argparse.ArgumentParser(description="Spot-check analytics.* tables")
    parser.add_argument("--database-url", default=os.environ.get("DATABASE_URL", ""))
    args = parser.parse_args()

    if not args.database_url.strip():
        print("DATABASE_URL is required", file=sys.stderr)
        return 2

    tables_sql = """
        SELECT table_name
        FROM information_schema.tables
        WHERE table_schema = 'analytics'
          AND table_type = 'BASE TABLE'
          AND table_name NOT LIKE 'sync_%'
        ORDER BY table_name
    """
    raw = run_psql(args.database_url, tables_sql)
    tables = [line.strip() for line in raw.splitlines() if line.strip()]
    if not tables:
        print("no analytics tables found", file=sys.stderr)
        return 2

    print(f"{'table':40} {'rows':>10} {'min_date':>12} {'max_date':>12}")
    print("-" * 78)

    date_columns = ("data", "refdate", "refmonth", "mes", "data_preco", "capturado_em")
    errors: list[str] = []

    for table in tables:
        count = run_psql(args.database_url, f"SELECT COUNT(*) FROM analytics.{table}")
        min_date = ""
        max_date = ""
        for col in date_columns:
            exists = run_psql(
                args.database_url,
                f"SELECT COUNT(*) FROM information_schema.columns "
                f"WHERE table_schema='analytics' AND table_name='{table}' AND column_name='{col}'",
            )
            if exists == "1":
                min_date = run_psql(args.database_url, f"SELECT MIN({col})::text FROM analytics.{table}")
                max_date = run_psql(args.database_url, f"SELECT MAX({col})::text FROM analytics.{table}")
                break
        print(f"{table:40} {count:>10} {min_date:>12} {max_date:>12}")
        if count == "0":
            errors.append(f"{table}: empty")

    if errors:
        print("\nspot-check failures:", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print(f"\nspot-check ok ({len(tables)} tables)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
