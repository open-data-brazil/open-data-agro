#!/usr/bin/env python3
"""Spot-check analytics tables for row counts and date ranges (PostgreSQL or DuckDB)."""

from __future__ import annotations

import argparse
import os
import subprocess
import sys

WAVE3_TABLES: tuple[str, ...] = (
    "dnit_snv_rodovias_federais",
    "ipea_series_macro_regionais",
    "ibge_pevs_producao_vegetal",
    "ibge_ppm_producao_municipal",
    "aneel_tarifas_energia",
    "bndes_financiamento_agro",
    "inmet_sequia_monitor",
    "oecd_ag_outlook",
    "fao_food_price_index",
    "argentina_magyp_producion_granos",
)

DATE_COLUMNS: tuple[str, ...] = (
    "data",
    "refdate",
    "refmonth",
    "refyear",
    "mes",
    "ano",
    "year",
    "data_preco",
    "capturado_em",
    "DatCompetencia",
)


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


def run_duckdb(duckdb_bin: str, duckdb_path: str, sql: str) -> str:
    proc = subprocess.run(
        [duckdb_bin, duckdb_path, "-csv", "-noheader", "-c", sql],
        check=False,
        capture_output=True,
        text=True,
    )
    if proc.returncode != 0:
        raise RuntimeError(proc.stderr.strip() or proc.stdout.strip())
    return proc.stdout.strip()


def list_postgres_tables(database_url: str) -> list[str]:
    tables_sql = """
        SELECT table_name
        FROM information_schema.tables
        WHERE table_schema = 'analytics'
          AND table_type = 'BASE TABLE'
          AND table_name NOT LIKE 'sync_%'
        ORDER BY table_name
    """
    raw = run_psql(database_url, tables_sql)
    return [line.strip() for line in raw.splitlines() if line.strip()]


def list_duckdb_tables(duckdb_bin: str, duckdb_path: str) -> list[str]:
    raw = run_duckdb(
        duckdb_bin,
        duckdb_path,
        "SELECT table_name FROM information_schema.tables "
        "WHERE table_schema = 'analytics' ORDER BY table_name",
    )
    return [line.strip() for line in raw.splitlines() if line.strip()]


def postgres_column_exists(database_url: str, table: str, column: str) -> bool:
    exists = run_psql(
        database_url,
        f"SELECT COUNT(*) FROM information_schema.columns "
        f"WHERE table_schema='analytics' AND table_name='{table}' AND column_name='{column}'",
    )
    return exists == "1"


def duckdb_column_exists(duckdb_bin: str, duckdb_path: str, table: str, column: str) -> bool:
    exists = run_duckdb(
        duckdb_bin,
        duckdb_path,
        f"SELECT COUNT(*) FROM information_schema.columns "
        f"WHERE table_schema='analytics' AND table_name='{table}' AND column_name='{column}'",
    )
    return exists == "1"


def spot_check_postgres(database_url: str, tables: list[str]) -> int:
    print(f"{'table':40} {'rows':>10} {'min_date':>12} {'max_date':>12}")
    print("-" * 78)
    errors: list[str] = []

    for table in tables:
        count = run_psql(database_url, f"SELECT COUNT(*) FROM analytics.{table}")
        min_date = ""
        max_date = ""
        for col in DATE_COLUMNS:
            if postgres_column_exists(database_url, table, col):
                min_date = run_psql(database_url, f"SELECT MIN({col})::text FROM analytics.{table}")
                max_date = run_psql(database_url, f"SELECT MAX({col})::text FROM analytics.{table}")
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


def spot_check_duckdb(duckdb_bin: str, duckdb_path: str, tables: list[str]) -> int:
    print(f"{'table':40} {'rows':>10} {'min_date':>12} {'max_date':>12}")
    print("-" * 78)
    errors: list[str] = []

    for table in tables:
        count = run_duckdb(duckdb_bin, duckdb_path, f"SELECT COUNT(*) FROM analytics.{table}")
        min_date = ""
        max_date = ""
        for col in DATE_COLUMNS:
            if duckdb_column_exists(duckdb_bin, duckdb_path, table, col):
                min_date = run_duckdb(
                    duckdb_bin, duckdb_path, f"SELECT MIN({col})::VARCHAR FROM analytics.{table}"
                )
                max_date = run_duckdb(
                    duckdb_bin, duckdb_path, f"SELECT MAX({col})::VARCHAR FROM analytics.{table}"
                )
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


def main() -> int:
    parser = argparse.ArgumentParser(description="Spot-check analytics.* tables")
    parser.add_argument("--database-url", default=os.environ.get("DATABASE_URL", ""))
    parser.add_argument("--duckdb", default=os.environ.get("DUCKDB_PATH", ""))
    parser.add_argument("--duckdb-bin", default=os.environ.get("DUCKDB_BIN", "duckdb"))
    parser.add_argument(
        "--wave3",
        action="store_true",
        help="check only wave 3 ingest marts (10 tables)",
    )
    args = parser.parse_args()

    use_duckdb = bool(args.duckdb.strip())
    use_postgres = bool(args.database_url.strip())
    if use_duckdb == use_postgres:
        print("exactly one of --duckdb or --database-url is required", file=sys.stderr)
        return 2

    if use_duckdb:
        tables = list_duckdb_tables(args.duckdb_bin, args.duckdb)
    else:
        tables = list_postgres_tables(args.database_url)

    if args.wave3:
        tables = [t for t in WAVE3_TABLES if t in tables]
        missing = [t for t in WAVE3_TABLES if t not in tables]
        if missing:
            print("missing wave 3 tables:", ", ".join(missing), file=sys.stderr)
            return 2

    if not tables:
        print("no analytics tables found", file=sys.stderr)
        return 2

    if use_duckdb:
        return spot_check_duckdb(args.duckdb_bin, args.duckdb, tables)
    return spot_check_postgres(args.database_url, tables)


if __name__ == "__main__":
    sys.exit(main())
