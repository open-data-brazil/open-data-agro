#!/usr/bin/env python3
"""Export soy daily ML training dataset with forward targets + walk-forward split (Phase 30)."""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import subprocess
import sys
from datetime import date, datetime, timezone
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.config import (  # noqa: E402
    FEATURE_COLS,
    HORIZONS,
    META_COLS,
    SCHEMA_VERSION,
    SPLIT_TEST_START,
    SPLIT_TRAIN_END,
    SPLIT_VAL_END,
    SPLIT_VAL_START,
)


def _git_sha() -> str:
    try:
        out = subprocess.check_output(
            ["git", "rev-parse", "HEAD"],
            cwd=ROOT,
            stderr=subprocess.DEVNULL,
            text=True,
        )
        return out.strip()
    except (subprocess.CalledProcessError, FileNotFoundError):
        return "unknown"


def _assign_split(d: date) -> str:
    if d <= date.fromisoformat(SPLIT_TRAIN_END):
        return "train"
    if date.fromisoformat(SPLIT_VAL_START) <= d <= date.fromisoformat(SPLIT_VAL_END):
        return "val"
    if d >= date.fromisoformat(SPLIT_TEST_START):
        return "test"
    return "holdout"


def _schema_hash(column_names: list[str]) -> str:
    payload = "|".join(column_names).encode("utf-8")
    return hashlib.sha256(payload).hexdigest()[:16]


def load_features(path: Path) -> pa.Table:
    if not path.is_file():
        raise FileNotFoundError(f"missing feature mart: {path}")
    return pq.read_table(path)


def build_training_table(features: pa.Table) -> pa.Table:
    import pyarrow.compute as pc

    required = {"data", "preco_rs_sc", "produto_slug"}
    missing = required - set(features.column_names)
    if missing:
        raise ValueError(f"feature mart missing columns: {sorted(missing)}")

    # Sort by date for lead windows.
    sort_idx = pc.sort_indices(features, sort_keys=[("data", "ascending")])
    features = features.take(sort_idx)

    prices = features.column("preco_rs_sc").combine_chunks()
    dates = features.column("data").combine_chunks()
    n = len(features)

    # Materialize to python for lead indexing (small panels / CI).
    price_list = [None if v is None else float(v) for v in prices.to_pylist()]
    date_list = dates.to_pylist()

    arrays: dict[str, list[object]] = {name: features.column(name).to_pylist() for name in features.column_names}

    for h in HORIZONS:
        ret_col: list[float | None] = []
        level_col: list[float | None] = []
        for i in range(n):
            j = i + h
            if j >= n or price_list[i] is None or price_list[j] is None or price_list[i] == 0:
                ret_col.append(None)
                level_col.append(None)
            else:
                level_col.append(price_list[j])
                ret_col.append(price_list[j] / price_list[i] - 1.0)
        arrays[f"y_ret_{h}d"] = ret_col
        arrays[f"y_level_{h}d"] = level_col

    splits: list[str] = []
    for d in date_list:
        if d is None:
            splits.append("holdout")
        else:
            dd = d if isinstance(d, date) else date.fromisoformat(str(d)[:10])
            splits.append(_assign_split(dd))
    arrays["split"] = splits

    # Keep feature + meta + targets only (stable schema).
    target_cols = [f"y_ret_{h}d" for h in HORIZONS] + [f"y_level_{h}d" for h in HORIZONS]
    ordered = [c for c in list(META_COLS) + list(FEATURE_COLS) + target_cols + ["split"] if c in arrays]
    # Deduplicate while preserving order
    seen: set[str] = set()
    final_cols: list[str] = []
    for c in ordered:
        if c not in seen:
            seen.add(c)
            final_cols.append(c)
    # Include any leftover source columns not listed (forward compatible)
    for c in arrays:
        if c not in seen:
            final_cols.append(c)

    return pa.table({c: arrays[c] for c in final_cols})


def write_export(table: pa.Table, out_dir: Path, source_path: Path) -> Path:
    out_dir.mkdir(parents=True, exist_ok=True)
    parquet_path = out_dir / "soy_daily_training.parquet"
    pq.write_table(table, parquet_path)

    dates = table.column("data").to_pylist()
    date_strs = [str(d)[:10] for d in dates if d is not None]
    splits = table.column("split").to_pylist()
    split_counts = {s: splits.count(s) for s in sorted(set(splits))}

    cols = table.column_names
    manifest = {
        "schema_version": SCHEMA_VERSION,
        "dataset_id": "ml.soy-daily-training",
        "produto_slug": "soja",
        "grain": ["produto_slug", "data"],
        "horizons_trading_days": list(HORIZONS),
        "target_definition": {
            "y_ret_hd": "preco_rs_sc[t+h] / preco_rs_sc[t] - 1 (trading-day lead h)",
            "y_level_hd": "preco_rs_sc[t+h] (trading-day lead h)",
        },
        "feature_cols": list(FEATURE_COLS),
        "split_policy": {
            "type": "walk_forward_calendar",
            "train": f"data <= {SPLIT_TRAIN_END}",
            "val": f"{SPLIT_VAL_START} .. {SPLIT_VAL_END}",
            "test": f"data >= {SPLIT_TEST_START}",
            "never": "random_shuffle",
        },
        "source_mart": str(source_path),
        "parquet": parquet_path.name,
        "row_count": table.num_rows,
        "column_count": table.num_columns,
        "columns": cols,
        "schema_hash": _schema_hash(cols),
        "date_min": min(date_strs) if date_strs else None,
        "date_max": max(date_strs) if date_strs else None,
        "split_counts": split_counts,
        "git_sha": _git_sha(),
        "created_at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "license_note": "CEPEA-derived columns are CC BY-NC 4.0 — non-commercial without separate license",
    }
    manifest_path = out_dir / "manifest.json"
    manifest_path.write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    return manifest_path


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--lake-root",
        default=os.environ.get("LAKE_LOCAL_ROOT", str(ROOT / "lake")),
    )
    parser.add_argument(
        "--out-dir",
        default=os.environ.get("ML_EXPORT_DIR", str(ROOT / ".local" / "ml" / "export")),
    )
    args = parser.parse_args()

    source = Path(args.lake_root) / "gold" / "mart_ml__soy_daily_features" / "mart.parquet"
    features = load_features(source)
    training = build_training_table(features)
    out_dir = Path(args.out_dir)
    manifest_path = write_export(training, out_dir, source)
    print(f"wrote {out_dir / 'soy_daily_training.parquet'} rows={training.num_rows}")
    print(f"wrote {manifest_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
