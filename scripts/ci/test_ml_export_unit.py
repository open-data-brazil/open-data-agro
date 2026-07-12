#!/usr/bin/env python3
"""Unit tests for Phase 30 export helpers (no network)."""

from __future__ import annotations

import sys
import tempfile
from datetime import date, timedelta
from pathlib import Path

import pyarrow as pa

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.export_soy_daily_dataset import build_training_table, write_export  # noqa: E402
from ml.verify_manifest import main as verify_main  # noqa: E402


def test_forward_targets_and_split() -> None:
    n = 40
    days = [date(2018, 1, 2) + timedelta(days=i) for i in range(n)]
    prices = [100.0 + i for i in range(n)]
    table = pa.table(
        {
            "produto_slug": ["soja"] * n,
            "data": days,
            "label_date": days,
            "preco_rs_sc": prices,
            "preco_usd_sc": prices,
            "variacao_dia_pct": [0.0] * n,
            "cepea_ret_1d": [0.0] * n,
            "cepea_ret_7d": [0.0] * n,
            "cepea_vol_7d": [0.0] * n,
            "ptax_usd_venda": [5.0] * n,
            "b3_front_price": [12.0] * n,
            "b3_front_symbol": ["X"] * n,
            "basis_cepea_usd_minus_b3": [0.0] * n,
            "comex_soja_fob_usd": [1.0] * n,
            "comex_soja_kg": [1.0] * n,
            "frete_mt_pr_ton_avg": [1.0] * n,
            "frete_mt_pr_tkm_avg": [1.0] * n,
            "inmet_national_avg": [1.0] * n,
            "as_of_cepea": days,
            "as_of_ptax": days,
            "as_of_b3": days,
            "as_of_comex": days,
            "as_of_frete": days,
            "as_of_inmet": days,
            "feature_as_of_max": days,
        }
    )
    out = build_training_table(table)
    assert "y_ret_7d" in out.column_names
    assert "split" in out.column_names
    assert abs(out.column("y_ret_7d")[0].as_py() - 0.07) < 1e-9
    assert out.column("y_level_7d")[0].as_py() == 107.0
    assert out.column("split")[0].as_py() == "train"


def test_manifest_roundtrip() -> None:
    n = 20
    days = [date(2022, 1, 3) + timedelta(days=i) for i in range(n)]
    prices = [100.0 + i for i in range(n)]
    table = pa.table(
        {
            "produto_slug": ["soja"] * n,
            "data": days,
            "label_date": days,
            "preco_rs_sc": prices,
            "preco_usd_sc": prices,
            "variacao_dia_pct": [0.0] * n,
            "cepea_ret_1d": [0.0] * n,
            "cepea_ret_7d": [0.0] * n,
            "cepea_vol_7d": [0.0] * n,
            "ptax_usd_venda": [5.0] * n,
            "b3_front_price": [12.0] * n,
            "b3_front_symbol": ["X"] * n,
            "basis_cepea_usd_minus_b3": [0.0] * n,
            "comex_soja_fob_usd": [1.0] * n,
            "comex_soja_kg": [1.0] * n,
            "frete_mt_pr_ton_avg": [1.0] * n,
            "frete_mt_pr_tkm_avg": [1.0] * n,
            "inmet_national_avg": [1.0] * n,
            "as_of_cepea": days,
            "as_of_ptax": days,
            "as_of_b3": days,
            "as_of_comex": days,
            "as_of_frete": days,
            "as_of_inmet": days,
            "feature_as_of_max": days,
        }
    )
    training = build_training_table(table)
    with tempfile.TemporaryDirectory() as tmp:
        out = Path(tmp)
        write_export(training, out, Path("fake/mart.parquet"))
        assert (out / "manifest.json").is_file()
        assert (out / "soy_daily_training.parquet").is_file()
        sys.argv = ["verify_manifest.py", "--export-dir", str(out)]
        assert verify_main() == 0


if __name__ == "__main__":
    test_forward_targets_and_split()
    test_manifest_roundtrip()
    print("test_ml_export_unit: PASS")
