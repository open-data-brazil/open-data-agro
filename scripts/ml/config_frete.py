"""Frozen frete forecast config (Phase 32 Track F)."""

from __future__ import annotations

# Primary target: R$/ton freight level.
TARGET_COL: str = "valor_frete_tonelada"
SECONDARY_TARGET_COL: str = "valor_tonelada_km"

# Calendar walk-forward on refmonth year (never random shuffle).
SPLIT_TRAIN_YEAR_END: int = 2022
SPLIT_VAL_YEAR_START: int = 2023
SPLIT_VAL_YEAR_END: int = 2023
SPLIT_TEST_YEAR_START: int = 2024

SUCCESS_IMPROVEMENT_MIN: float = 0.15  # beat seasonal_naive MAE by ≥15% on test
SUCCESS_BASELINE: str = "seasonal_naive"

# Numeric + categorical model inputs (categoricals passed as codes).
NUMERIC_FEATURE_COLS: tuple[str, ...] = (
    "distancia_km",
    "mes",
    "ano",
    "corridor_lag12_ton_avg",
    "corridor_lag1_ton_avg",
)
CATEGORICAL_FEATURE_COLS: tuple[str, ...] = (
    "uf_origem",
    "uf_destino",
    "fonte_code",
)

META_COLS: tuple[str, ...] = (
    "municipio_origem",
    "municipio_destino",
    "cod_ibge_origem",
    "cod_ibge_destino",
    "refmonth",
    "as_of",
    "split",
    "fonte",
)

SCHEMA_VERSION: str = "frete_corridor_v1"
MIN_EXPORT_ROWS: int = 500
