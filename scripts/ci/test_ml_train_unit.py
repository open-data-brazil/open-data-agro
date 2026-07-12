#!/usr/bin/env python3
"""Unit tests for Phase 31 train helpers (no network)."""

from __future__ import annotations

import math
import sys
import tempfile
from datetime import date, timedelta
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.config import FEATURE_COLS, SUCCESS_HORIZON  # noqa: E402
from ml.metrics import gate_passes, mae  # noqa: E402
from ml.train_soy import load_rows, split_xy, train_horizon, write_registry  # noqa: E402


def test_gate_passes() -> None:
    assert gate_passes(0.85, 1.0, 0.15) is True
    assert gate_passes(0.86, 1.0, 0.15) is False
    assert gate_passes(0.1, 0.0, 0.15) is False


def test_mae() -> None:
    assert mae([1.0, 2.0], [1.0, 3.0]) == 0.5
    assert mae([], []) is None


def _business_days(start: date, count: int) -> list[date]:
    days: list[date] = []
    d = start
    while len(days) < count:
        if d.weekday() < 5:
            days.append(d)
        d += timedelta(days=1)
    return days


def _synthetic_rows(n: int = 500) -> list[dict[str, object]]:
    days = _business_days(date(2016, 1, 4), n)
    prices = [100.0]
    for i in range(1, n):
        driver = 0.03 * math.sin(i / 18.0)
        prices.append(prices[-1] * (1.0 + 0.001 * driver))

    rows: list[dict[str, object]] = []
    for i, d in enumerate(days):
        driver = 0.03 * math.sin(i / 18.0)
        if d.year <= 2018:
            split = "train"
        elif d.year <= 2021:
            split = "val"
        else:
            split = "test"
        row: dict[str, object] = {
            "data": d.isoformat(),
            "split": split,
            "preco_rs_sc": prices[i],
        }
        for col in FEATURE_COLS:
            if col == "preco_rs_sc":
                continue
            if col in ("b3_front_price", "basis_cepea_usd_minus_b3", "frete_mt_pr_ton_avg"):
                row[col] = 10.0 * driver
            else:
                row[col] = float(i % 7)
        rows.append(row)

    for h in (7, 30, 90):
        for i, row in enumerate(rows):
            j = i + h
            if j < len(rows):
                p0 = float(row["preco_rs_sc"])  # type: ignore[arg-type]
                p1 = float(rows[j]["preco_rs_sc"])  # type: ignore[arg-type]
                row[f"y_ret_{h}d"] = p1 / p0 - 1.0
                row[f"y_level_{h}d"] = p1
            else:
                row[f"y_ret_{h}d"] = None
                row[f"y_level_{h}d"] = None
    return rows


def test_split_xy_and_train_smoke() -> None:
    try:
        import lightgbm  # noqa: F401
    except ImportError:
        print("skip train smoke: lightgbm not installed")
        return

    # ~9y business days so train/val/test splits are non-empty.
    rows = _synthetic_rows(2300)
    target = f"y_ret_{SUCCESS_HORIZON}d"
    x_train, y_train = split_xy(rows, target, "train")
    assert len(y_train) > 50
    assert x_train.shape[1] == len(FEATURE_COLS)

    from ml.train_soy import _require_lightgbm

    lgb = _require_lightgbm()
    result = train_horizon(lgb, rows, SUCCESS_HORIZON)
    assert result["n_test"] > 0  # type: ignore[index]
    assert result["point"]["test_mae"] is not None  # type: ignore[index]
    assert "quantiles" in result
    assert set(result["quantiles"]) >= {"p10", "p50", "p90"}  # type: ignore[arg-type]


def test_registry_write() -> None:
    with tempfile.TemporaryDirectory() as tmp:
        path = write_registry(Path(tmp), {"b3_gate_passed": True}, "abc123")
        assert path.is_file()
        text = path.read_text(encoding="utf-8")
        assert "soy_lgbm_walkforward_v1" in text
        assert "abc123" in text


def test_load_rows_roundtrip() -> None:
    rows = _synthetic_rows(30)
    table = pa.Table.from_pylist(rows)
    with tempfile.TemporaryDirectory() as tmp:
        path = Path(tmp) / "soy_daily_training.parquet"
        pq.write_table(table, path)
        loaded = load_rows(path)
        assert len(loaded) == 30


def main() -> int:
    test_gate_passes()
    test_mae()
    test_registry_write()
    test_load_rows_roundtrip()
    test_split_xy_and_train_smoke()
    print("test_ml_train_unit: PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
