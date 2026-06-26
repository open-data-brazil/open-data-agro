#!/usr/bin/env python3
"""Validate CONAB IBGE municipality codes against the IBGE localidades mart."""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

import pandas as pd

REPO_ROOT = Path(__file__).resolve().parents[2]

# Sentinel codes used by CONAB when origin municipality is unknown.
SENTINEL_CODES = frozenset({"9999999"})

# CONAB gold marts (relative to lake/gold) and IBGE code columns to check.
CONAB_IBGE_COLUMNS: dict[str, list[str]] = {
    "mart_conab__precos_semanal_municipio": ["cod_ibge"],
    "mart_conab__precos_mensal_municipio": ["cod_ibge"],
    "mart_conab__frete": ["cod_ibge_origem", "cod_ibge_destino"],
    "mart_conab__armazenagem": ["cod_ibge"],
    "mart_conab__estoques_publicos": ["cod_ibge"],
    "mart_conab__alimenta_brasil_propostas": ["cod_ibge"],
    "mart_conab__prohort_diario": ["cod_ibge_municipio"],
    "mart_conab__prohort_mensal": [
        "cod_ibge_municipio_ceasa",
        "cod_ibge_municipio_origem",
    ],
}

IBGE_MUNICIPIOS_MART = "mart_ibge__localidades_municipios"


def normalize_codigo_ibge(value: object) -> str | None:
    if value is None or (isinstance(value, float) and pd.isna(value)):
        return None
    text = str(value).strip()
    if not text:
        return None
    if text.isdigit():
        return text.zfill(7)
    return text


def load_municipio_codes(lake_root: Path) -> set[str]:
    mart_path = lake_root / "gold" / IBGE_MUNICIPIOS_MART / "mart.parquet"
    if not mart_path.is_file():
        raise FileNotFoundError(
            f"IBGE municipios mart not found: {mart_path} "
            "(run make ibge-localidades-mvp or ingest ibge.localidades-municipios)"
        )
    df = pd.read_parquet(mart_path, columns=["codigo_ibge"])
    codes = {normalize_codigo_ibge(v) for v in df["codigo_ibge"]}
    return {code for code in codes if code}


def validate_mart(
    mart_path: Path,
    columns: list[str],
    reference: set[str],
    *,
    sample_limit: int,
) -> dict:
    df = pd.read_parquet(mart_path, columns=columns)
    result: dict = {
        "mart": mart_path.parent.name,
        "path": str(mart_path),
        "columns": {},
        "ok": True,
    }

    for column in columns:
        if column not in df.columns:
            result["columns"][column] = {
                "ok": False,
                "error": f"missing column {column}",
            }
            result["ok"] = False
            continue

        invalid: dict[str, int] = {}
        bad_length: dict[str, int] = {}
        checked = 0

        for raw in df[column].dropna().unique():
            code = normalize_codigo_ibge(raw)
            if code is None:
                continue
            checked += 1
            if len(code) != 7 or not code.isdigit():
                bad_length[code] = int((df[column].astype(str).str.strip() == str(raw).strip()).sum())
                continue
            if code in SENTINEL_CODES:
                continue
            if code not in reference:
                invalid[code] = int((df[column].map(normalize_codigo_ibge) == code).sum())

        column_ok = not invalid and not bad_length
        sample_invalid = dict(sorted(invalid.items(), key=lambda item: -item[1])[:sample_limit])
        sample_bad_length = dict(sorted(bad_length.items(), key=lambda item: -item[1])[:sample_limit])

        result["columns"][column] = {
            "ok": column_ok,
            "distinct_checked": checked,
            "invalid_count": len(invalid),
            "invalid_rows": int(sum(invalid.values())),
            "invalid_sample": sample_invalid,
            "bad_length_count": len(bad_length),
            "bad_length_sample": sample_bad_length,
        }
        if not column_ok:
            result["ok"] = False

    return result


def run_validation(lake_root: Path, *, sample_limit: int = 10) -> dict:
    reference = load_municipio_codes(lake_root)
    reports: list[dict] = []
    skipped: list[str] = []

    for mart_name, columns in CONAB_IBGE_COLUMNS.items():
        mart_path = lake_root / "gold" / mart_name / "mart.parquet"
        if not mart_path.is_file():
            skipped.append(mart_name)
            continue
        reports.append(
            validate_mart(mart_path, columns, reference, sample_limit=sample_limit)
        )

    return {
        "lake_root": str(lake_root),
        "reference_municipios": len(reference),
        "validated_marts": len(reports),
        "skipped_marts": skipped,
        "ok": all(report["ok"] for report in reports),
        "reports": reports,
    }


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Validate CONAB cod_ibge values against IBGE municipios mart"
    )
    parser.add_argument(
        "--lake-root",
        default=str(REPO_ROOT / "lake"),
        help="Lake root containing gold/ marts (default: ./lake)",
    )
    parser.add_argument(
        "--sample-limit",
        type=int,
        default=10,
        help="Max invalid codes to show per column (default: 10)",
    )
    parser.add_argument(
        "--json",
        action="store_true",
        help="Print machine-readable JSON report",
    )
    args = parser.parse_args()

    lake_root = Path(args.lake_root).resolve()
    try:
        summary = run_validation(lake_root, sample_limit=args.sample_limit)
    except FileNotFoundError as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 2

    if args.json:
        print(json.dumps(summary, indent=2, ensure_ascii=False))
    else:
        print(f"lake={summary['lake_root']}")
        print(f"reference_municipios={summary['reference_municipios']}")
        print(f"validated_marts={summary['validated_marts']}")
        if summary["skipped_marts"]:
            print(f"skipped_marts={', '.join(summary['skipped_marts'])}")
        for report in summary["reports"]:
            status = "OK" if report["ok"] else "FAIL"
            print(f"\n[{status}] {report['mart']}")
            for column, detail in report["columns"].items():
                col_status = "OK" if detail.get("ok") else "FAIL"
                print(f"  {column}: {col_status}")
                if detail.get("invalid_sample"):
                    print(f"    invalid_sample={detail['invalid_sample']}")
                if detail.get("bad_length_sample"):
                    print(f"    bad_length_sample={detail['bad_length_sample']}")
        print(f"\noverall={'PASS' if summary['ok'] else 'FAIL'}")

    return 0 if summary["ok"] else 1


if __name__ == "__main__":
    sys.exit(main())
