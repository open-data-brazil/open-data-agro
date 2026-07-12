"""Shared regression / gate metrics for Phase 30–31."""

from __future__ import annotations

import math


def mae(y_true: list[float], y_pred: list[float]) -> float | None:
    if not y_true:
        return None
    return sum(abs(a - b) for a, b in zip(y_true, y_pred, strict=True)) / len(y_true)


def rmse(y_true: list[float], y_pred: list[float]) -> float | None:
    if not y_true:
        return None
    return math.sqrt(sum((a - b) ** 2 for a, b in zip(y_true, y_pred, strict=True)) / len(y_true))


def gate_passes(model_mae: float, baseline_mae: float, min_improvement: float) -> bool:
    """True when model MAE beats baseline by ≥ min_improvement (relative)."""
    if baseline_mae <= 0:
        return False
    return model_mae <= (1.0 - min_improvement) * baseline_mae
