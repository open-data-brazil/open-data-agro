"""Frozen Phase 30 training-dataset config (UC-ML-001)."""

from __future__ import annotations

# Trading-day horizons: lead by h CEPEA observations (not calendar days).
HORIZONS: tuple[int, ...] = (7, 30, 90)

# Primary success-gate horizon (B3).
SUCCESS_HORIZON: int = 30
SUCCESS_METRIC: str = "mae"
SUCCESS_BASELINE: str = "seasonal_naive"
SUCCESS_IMPROVEMENT_MIN: float = 0.15  # beat baseline MAE by ≥15% on test OOS

FEATURE_COLS: tuple[str, ...] = (
    "preco_rs_sc",
    "preco_usd_sc",
    "variacao_dia_pct",
    "cepea_ret_1d",
    "cepea_ret_7d",
    "cepea_vol_7d",
    "ptax_usd_venda",
    "b3_front_price",
    "basis_cepea_usd_minus_b3",
    "comex_soja_fob_usd",
    "comex_soja_kg",
    "frete_mt_pr_ton_avg",
    "frete_mt_pr_tkm_avg",
    "inmet_national_avg",
)

# Columns excluded from model features (ids / leakage audit / symbols).
META_COLS: tuple[str, ...] = (
    "produto_slug",
    "data",
    "label_date",
    "b3_front_symbol",
    "as_of_cepea",
    "as_of_ptax",
    "as_of_b3",
    "as_of_comex",
    "as_of_frete",
    "as_of_inmet",
    "feature_as_of_max",
    "split",
)

# Walk-forward calendar splits on column `data` (never random shuffle).
SPLIT_TRAIN_END: str = "2018-12-31"
SPLIT_VAL_START: str = "2019-01-01"
SPLIT_VAL_END: str = "2021-12-31"
SPLIT_TEST_START: str = "2022-01-01"

# Lagged correlation feature candidates vs forward return target.
CORR_FEATURE_CANDIDATES: tuple[str, ...] = (
    "cepea_ret_1d",
    "cepea_ret_7d",
    "cepea_vol_7d",
    "ptax_usd_venda",
    "basis_cepea_usd_minus_b3",
    "comex_soja_kg",
    "frete_mt_pr_ton_avg",
    "inmet_national_avg",
)
CORR_LAGS: tuple[int, ...] = (0, 1, 7, 30)
CORR_TARGET: str = "y_ret_30d"

SCHEMA_VERSION: str = "soy_daily_v1"
