#!/usr/bin/env python3
"""Benchmark ingestor end-to-end: download, convert, bronze land, Postgres audit.

Usage:
  python3 scripts/benchmark/ingestor_benchmark.py --clean --all
  python3 scripts/benchmark/ingestor_benchmark.py --profile scripts/benchmark/profiles/fast10.json
  python3 scripts/benchmark/ingestor_benchmark.py --profile scripts/benchmark/profiles/fast10-stress.json
  python3 scripts/benchmark/ingestor_benchmark.py --datasets bcb.sgs-ipca,cepea.soja-paranagua
  make benchmark-ingestor-fast10-clean
  make benchmark-ingestor-fast10-stress

Requires: postgres up, migrations applied, .env with DATABASE_URL.
"""

from __future__ import annotations

import argparse
import json
import os
import shutil
import subprocess
import sys
import time
from collections import defaultdict
from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parents[2]

DEFAULT_DATASETS = [
    "ibge.localidades-ufs",
    "bcb.sgs-ipca",
    "bcb.sgs-ptax-usd-venda",
    "cepea.soja-paranagua",
    "conab.estimativa-graos",
]


@dataclass
class BenchmarkRow:
    dataset_id: str
    status: str
    duration_sec: float
    row_count: int | None
    file_size_bytes: int | None
    sha256: str
    bronze_key: str
    skipped: bool
    error: str


def load_dotenv(path: Path) -> None:
    if not path.is_file():
        return
    for line in path.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, _, value = line.partition("=")
        os.environ.setdefault(key.strip(), value.strip())


def psql_exec(sql: str) -> None:
    database_url = os.environ.get("DATABASE_URL", "")
    if not database_url:
        raise SystemExit("DATABASE_URL not set")

    docker_ps = subprocess.run(
        ["docker", "compose", "ps", "-q", "postgres"],
        cwd=REPO_ROOT,
        capture_output=True,
        text=True,
        check=False,
    )
    if docker_ps.stdout.strip():
        cmd = [
            "docker", "compose", "exec", "-T", "postgres",
            "psql", "-U", "open_data_agro", "-d", "open_data_agro", "-c", sql,
        ]
        proc = subprocess.run(cmd, cwd=REPO_ROOT, capture_output=True, text=True, check=False)
    else:
        proc = subprocess.run(
            ["psql", database_url, "-c", sql],
            capture_output=True,
            text=True,
            check=False,
        )
    if proc.returncode != 0:
        raise RuntimeError(proc.stderr.strip() or proc.stdout.strip())


def psql_query(sql: str) -> list[list[str]]:
    database_url = os.environ.get("DATABASE_URL", "")
    if not database_url:
        raise SystemExit("DATABASE_URL not set")

    docker_ps = subprocess.run(
        ["docker", "compose", "ps", "-q", "postgres"],
        cwd=REPO_ROOT,
        capture_output=True,
        text=True,
        check=False,
    )
    if docker_ps.stdout.strip():
        cmd = [
            "docker", "compose", "exec", "-T", "postgres",
            "psql", "-U", "open_data_agro", "-d", "open_data_agro",
            "-t", "-A", "-F", "\t", "-c", sql,
        ]
        proc = subprocess.run(cmd, cwd=REPO_ROOT, capture_output=True, text=True, check=False)
    else:
        proc = subprocess.run(
            ["psql", database_url, "-t", "-A", "-F", "\t", "-c", sql],
            capture_output=True,
            text=True,
            check=False,
        )
    if proc.returncode != 0:
        raise RuntimeError(proc.stderr.strip() or proc.stdout.strip())
    rows: list[list[str]] = []
    for line in proc.stdout.splitlines():
        if line.strip():
            rows.append(line.split("\t"))
    return rows


def clean_state() -> None:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "./lake")).resolve()
    bronze = lake_root / "bronze"
    if bronze.is_dir():
        for child in bronze.iterdir():
            if child.name == ".gitkeep":
                continue
            if child.is_dir():
                shutil.rmtree(child)
            else:
                child.unlink(missing_ok=True)

    psql_exec("TRUNCATE catalog.ingest_files, catalog.ingest_jobs, catalog.dataset_registry CASCADE;")
    print(f"Clean slate: bronze under {bronze} + catalog ingest tables truncated", flush=True)


def list_catalog_datasets() -> list[str]:
    proc = subprocess.run(
        ["go", "run", "./cmd/ingestor", "catalog", "list"],
        cwd=REPO_ROOT,
        capture_output=True,
        text=True,
        check=False,
    )
    if proc.returncode != 0:
        raise RuntimeError(proc.stderr.strip() or proc.stdout.strip())
    datasets: list[str] = []
    for line in proc.stdout.splitlines():
        dataset_id = line.split("\t", 1)[0].strip()
        if dataset_id:
            datasets.append(dataset_id)
    return datasets


def load_profile(path: Path) -> dict:
    payload = json.loads(path.read_text(encoding="utf-8"))
    if "datasets" not in payload:
        raise SystemExit(f"profile {path} missing 'datasets' list")
    return payload


def normalize_profile_datasets(raw_datasets: list) -> tuple[list[str], dict[str, dict]]:
    dataset_ids: list[str] = []
    overrides: dict[str, dict] = {}
    for item in raw_datasets:
        if isinstance(item, str):
            dataset_ids.append(item)
            continue
        if isinstance(item, dict) and item.get("id"):
            dataset_ids.append(str(item["id"]))
            args = item.get("args") or {}
            if args:
                overrides[str(item["id"])] = args
            continue
        raise SystemExit(f"invalid profile dataset entry: {item!r}")
    return dataset_ids, overrides


def profile_defaults(profile: dict) -> dict:
    defaults = profile.get("defaults") or {}
    return {
        "cepea_from": str(defaults.get("cepea_from", "2010-01-01")),
        "pam_from": int(defaults.get("pam_from", 2010)),
        "pam_to": int(defaults.get("pam_to", 2024)),
        "pam_crop": str(defaults.get("pam_crop", "")),
        "pam_uf": str(defaults.get("pam_uf", "")),
        "inmet_year": int(defaults.get("inmet_year", 2024)),
        "inmet_uf": str(defaults.get("inmet_uf", "")),
    }


def build_ingest_args(
    dataset_id: str,
    global_from: str,
    defaults: dict | None = None,
    overrides: dict[str, dict] | None = None,
) -> list[str]:
    """Catalog-aware CLI flags; profile defaults/overrides take precedence."""
    cfg = defaults or {}
    extra: list[str] = []
    ds_override = (overrides or {}).get(dataset_id, {})

    if dataset_id.startswith("cepea."):
        extra.extend(["--from", str(ds_override.get("from", global_from or cfg.get("cepea_from", "2010-01-01")))])
    elif dataset_id.startswith("ibge.pam-"):
        extra.extend([
            "--from", str(ds_override.get("from", cfg.get("pam_from", 2010))),
            "--to", str(ds_override.get("to", cfg.get("pam_to", 2024))),
        ])
        crop = ds_override.get("crop", cfg.get("pam_crop", ""))
        if crop:
            extra.extend(["--crop", str(crop)])
        uf = ds_override.get("uf", cfg.get("pam_uf", ""))
        if uf:
            extra.extend(["--uf", str(uf)])
    elif dataset_id in (
        "inmet.bdmep-diario",
        "inmet.bdmep-mensal",
        "inmet.pacote-anual-automaticas",
    ):
        extra.extend(["--year", str(ds_override.get("year", cfg.get("inmet_year", 2024)))])
        uf = ds_override.get("uf", cfg.get("inmet_uf", ""))
        if uf:
            extra.extend(["--uf", str(uf)])

    for key in ("from", "to", "crop", "uf", "year"):
        if key in ds_override and key not in _args_keys(extra):
            value = ds_override[key]
            flag = f"--{key}"
            if flag not in extra:
                extra.extend([flag, str(value)])

    return extra


def _args_keys(extra: list[str]) -> set[str]:
    keys: set[str] = set()
    for i, token in enumerate(extra):
        if token.startswith("--") and i + 1 < len(extra):
            keys.add(token.lstrip("-"))
    return keys


def ingest_extra_args(dataset_id: str, global_from: str) -> list[str]:
    """Backward-compatible wrapper for non-profile runs."""
    return build_ingest_args(dataset_id, global_from)


def latest_job_metrics(dataset_id: str) -> dict[str, str]:
    rows = psql_query(
        f"""
        SELECT j.status,
               COALESCE(EXTRACT(EPOCH FROM (j.finished_at - j.started_at))::text, ''),
               COALESCE(f.row_count::text, ''),
               COALESCE(f.file_size_bytes::text, ''),
               COALESCE(f.sha256, ''),
               COALESCE(f.r2_key, ''),
               COALESCE(j.error_message, '')
        FROM catalog.ingest_jobs j
        LEFT JOIN catalog.ingest_files f ON f.job_id = j.id
        WHERE j.dataset_id = '{dataset_id.replace("'", "''")}'
          AND j.dry_run = false
        ORDER BY j.started_at DESC
        LIMIT 1
        """
    )
    if not rows:
        return {}
    cols = rows[0]
    keys = ["status", "duration_db", "row_count", "file_size_bytes", "sha256", "bronze_key", "error"]
    return dict(zip(keys, (cols + [""] * len(keys))[: len(keys)]))


def bronze_size_on_disk(bronze_key: str) -> int | None:
    if not bronze_key:
        return None
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "./lake"))
    path = lake_root / bronze_key
    if path.is_file():
        return path.stat().st_size
    return None


def ingestor_cmd(dataset_id: str, extra: list[str]) -> list[str]:
    binary = REPO_ROOT / "bin" / "ingestor"
    if binary.is_file():
        return [str(binary), "run", dataset_id, *extra]
    return ["go", "run", "./cmd/ingestor", "run", dataset_id, *extra]


def run_ingest(dataset_id: str, extra: list[str]) -> tuple[str, float, str]:
    cmd = ingestor_cmd(dataset_id, extra)
    started = time.perf_counter()
    proc = subprocess.run(
        cmd,
        cwd=REPO_ROOT,
        capture_output=True,
        text=True,
        env=os.environ.copy(),
        check=False,
    )
    elapsed = time.perf_counter() - started
    output = (proc.stdout or "") + (proc.stderr or "")
    if proc.returncode != 0:
        return "failed", elapsed, output.strip()
    return "ok", elapsed, output.strip()


def parse_ingestor_line(output: str) -> dict[str, str]:
    line = output.strip().splitlines()[-1] if output else ""
    out: dict[str, str] = {}
    for part in line.split():
        if "=" in part:
            k, v = part.split("=", 1)
            out[k] = v
    return out


def benchmark_dataset(
    dataset_id: str,
    global_from: str,
    defaults: dict | None = None,
    overrides: dict[str, dict] | None = None,
) -> BenchmarkRow:
    extra = build_ingest_args(dataset_id, global_from, defaults, overrides)
    flags = " ".join(extra)
    print(f"→ {dataset_id} {flags}".rstrip(), flush=True)
    _, wall_sec, output = run_ingest(dataset_id, extra)
    parsed = parse_ingestor_line(output)
    status = parsed.get("status", "failed")
    skipped = "skipped" in output.lower() or status == "skipped"
    sha256 = parsed.get("sha256", "")

    if status == "failed" and output:
        err_tail = output.splitlines()[-1][:200]
        return BenchmarkRow(
            dataset_id=dataset_id,
            status="failed",
            duration_sec=round(wall_sec, 3),
            row_count=None,
            file_size_bytes=None,
            sha256=sha256,
            bronze_key="",
            skipped=False,
            error=err_tail,
        )

    metrics = latest_job_metrics(dataset_id)
    if metrics:
        status = metrics.get("status", status)
        if metrics.get("error"):
            status = "failed"

    row_count = int(metrics["row_count"]) if metrics.get("row_count") else None
    file_size = int(metrics["file_size_bytes"]) if metrics.get("file_size_bytes") else None
    bronze_key = metrics.get("bronze_key", parsed.get("key", ""))
    disk_size = bronze_size_on_disk(bronze_key)
    if disk_size is not None:
        file_size = disk_size

    db_duration = float(metrics["duration_db"]) if metrics.get("duration_db") else wall_sec

    return BenchmarkRow(
        dataset_id=dataset_id,
        status=status,
        duration_sec=round(db_duration if db_duration > 0 else wall_sec, 3),
        row_count=row_count,
        file_size_bytes=file_size,
        sha256=sha256 or metrics.get("sha256", ""),
        bronze_key=bronze_key,
        skipped=skipped,
        error=metrics.get("error", "") if status == "failed" else "",
    )


def format_bytes(n: int | None) -> str:
    if n is None:
        return "—"
    value = float(n)
    for unit in ("B", "KB", "MB", "GB"):
        if value < 1024:
            return f"{int(value)} B" if unit == "B" else f"{value:.1f} {unit}"
        value /= 1024
    return f"{value:.1f} TB"


def agency_of(dataset_id: str) -> str:
    return dataset_id.split(".", 1)[0]


def print_report(
    rows: list[BenchmarkRow],
    total_sec: float,
    clean: bool,
    profile: dict | None = None,
) -> None:
    success = [r for r in rows if r.status == "success" and not r.skipped]
    skipped = [r for r in rows if r.skipped]
    failed = [r for r in rows if r.status not in ("success", "skipped")]

    total_rows = sum(r.row_count or 0 for r in success)
    total_bytes = sum(r.file_size_bytes or 0 for r in success)
    ingest_time = sum(r.duration_sec for r in success)

    by_agency: dict[str, list[BenchmarkRow]] = defaultdict(list)
    for r in success:
        by_agency[agency_of(r.dataset_id)].append(r)

    title = "INGESTOR CLEAN BENCHMARK (from zero)" if clean else "INGESTOR BENCHMARK REPORT"
    if profile:
        title = f"{title} — profile {profile.get('name', 'custom')}"
    print()
    print("=" * 80)
    print(title)
    if profile and profile.get("description"):
        print(profile["description"])
    print(f"Generated: {datetime.now(timezone.utc).strftime('%Y-%m-%d %H:%M:%S UTC')}")
    print("=" * 80)
    print()
    print(f"{'Dataset':<40} {'Status':<8} {'Time':>8} {'Rows':>12} {'Bronze':>10}")
    print("-" * 80)
    for r in rows:
        rows_s = f"{r.row_count:,}" if r.row_count is not None else "—"
        size_s = format_bytes(r.file_size_bytes)
        flag = " (skip)" if r.skipped else ""
        print(f"{r.dataset_id:<40} {r.status:<8} {r.duration_sec:>7.2f}s {rows_s:>12} {size_s:>10}{flag}")
        if r.error:
            print(f"  error: {r.error[:160]}")

    print("-" * 80)
    print(
        f"{'TOTAL (success only; failures excluded)':<40} {'':8} "
        f"{ingest_time:>7.2f}s {total_rows:>12,} {format_bytes(total_bytes):>10}"
    )
    if failed:
        print()
        print("Excluded from totals (failed)")
        for r in failed:
            print(f"  - {r.dataset_id}: {r.error[:120] if r.error else r.status}")
    print()
    print("Summary")
    print(f"  Datasets run:       {len(rows)}")
    print(f"  Success:              {len(success)}")
    print(f"  Skipped:              {len(skipped)}")
    print(f"  Failed:               {len(failed)}")
    print(f"  Wall clock:           {total_sec:.2f}s")
    print(f"  Sum ingest time:      {ingest_time:.2f}s")
    print(f"  Total rows (bronze):  {total_rows:,}")
    print(f"  Total bronze size:    {format_bytes(total_bytes)}")
    if total_sec > 0 and total_rows > 0:
        print(f"  Throughput:           {total_rows / total_sec:,.0f} rows/s (wall clock)")
    if total_sec > 0 and total_bytes > 0:
        print(f"  Data rate:            {total_bytes / total_sec / 1024:.1f} KB/s (wall clock)")

    budget = int(profile.get("budget_wall_clock_sec", 0)) if profile else 0
    if budget > 0:
        status = "OK" if total_sec <= budget else "EXCEEDED"
        print(f"  Budget ({budget}s):       {status} ({total_sec:.1f}s wall clock)")

    if by_agency:
        print()
        print("By agency (success)")
        print(f"  {'Agency':<10} {'Datasets':>8} {'Rows':>14} {'Bronze':>12} {'Time':>10}")
        for agency in sorted(by_agency):
            items = by_agency[agency]
            a_rows = sum(x.row_count or 0 for x in items)
            a_bytes = sum(x.file_size_bytes or 0 for x in items)
            a_time = sum(x.duration_sec for x in items)
            print(f"  {agency:<10} {len(items):>8} {a_rows:>14,} {format_bytes(a_bytes):>12} {a_time:>9.1f}s")

    print()
    print(f"Bronze lake:  {Path(os.environ.get('LAKE_LOCAL_ROOT', './lake')).resolve()}")
    print(f"Postgres:     catalog.ingest_jobs / catalog.ingest_files")
    print("=" * 80)


def main() -> int:
    parser = argparse.ArgumentParser(description="Benchmark ingestor datasets")
    parser.add_argument("--all", action="store_true", help="Run every dataset in the catalog")
    parser.add_argument("--clean", action="store_true", help="Truncate ingest audit + wipe bronze before run")
    parser.add_argument("--datasets", default="", help="Comma-separated dataset IDs")
    parser.add_argument("--profile", default="", help="JSON profile path (e.g. scripts/benchmark/profiles/fast10.json)")
    parser.add_argument("--from", dest="from_date", default="", help="Global --from for CEPEA datasets")
    parser.add_argument("--json", dest="json_out", default="", help="Write JSON report to path")
    args = parser.parse_args()

    load_dotenv(REPO_ROOT / ".env")

    profile: dict | None = None
    defaults: dict | None = None
    overrides: dict[str, dict] = {}

    if args.profile:
        profile_path = Path(args.profile)
        if not profile_path.is_absolute():
            profile_path = REPO_ROOT / profile_path
        profile = load_profile(profile_path)
        datasets, overrides = normalize_profile_datasets(profile["datasets"])
        defaults = profile_defaults(profile)
    elif args.all:
        datasets = list_catalog_datasets()
    elif args.datasets:
        datasets = [d.strip() for d in args.datasets.split(",") if d.strip()]
    else:
        datasets = DEFAULT_DATASETS

    if args.clean:
        clean_state()

    label = profile.get("name", "custom") if profile else f"{len(datasets)} datasets"
    print(f"Benchmarking {label} ({len(datasets)} datasets)...", flush=True)

    started = time.perf_counter()
    rows = [
        benchmark_dataset(d, args.from_date, defaults, overrides)
        for d in datasets
    ]
    total_sec = time.perf_counter() - started

    print_report(rows, total_sec, clean=args.clean, profile=profile)

    if args.json_out:
        json_path = args.json_out
    elif profile:
        suffix = "-clean" if args.clean else ""
        json_path = f".local/benchmark/ingestor-{profile.get('name', 'profile')}{suffix}.json"
    elif args.clean and args.all:
        json_path = ".local/benchmark/ingestor-clean-full.json"
    else:
        json_path = ""
    if json_path:
        success = [r for r in rows if r.status == "success" and not r.skipped]
        payload = {
            "generated_at": datetime.now(timezone.utc).isoformat(),
            "profile": profile.get("name") if profile else None,
            "clean_run": args.clean,
            "budget_wall_clock_sec": profile.get("budget_wall_clock_sec") if profile else None,
            "budget_exceeded": bool(profile and total_sec > int(profile.get("budget_wall_clock_sec", 0))),
            "dataset_count": len(rows),
            "success_count": len(success),
            "failed_count": len([r for r in rows if r.status not in ("success", "skipped")]),
            "excluded_from_totals": [r.dataset_id for r in rows if r.status not in ("success", "skipped")],
            "total_duration_sec": round(total_sec, 3),
            "total_rows": sum(r.row_count or 0 for r in success),
            "total_bronze_bytes": sum(r.file_size_bytes or 0 for r in success),
            "rows": [asdict(r) for r in rows],
        }
        out = REPO_ROOT / json_path
        out.parent.mkdir(parents=True, exist_ok=True)
        out.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
        print(f"JSON report: {out}")

    failed = [r for r in rows if r.status not in ("success", "skipped")]
    return 1 if failed else 0


if __name__ == "__main__":
    sys.exit(main())
