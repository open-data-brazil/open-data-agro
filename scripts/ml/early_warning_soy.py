#!/usr/bin/env python3
"""Early-warning rules + precision/recall on historical spikes (Phase 31 C3)."""

from __future__ import annotations

import argparse
import json
import os
import sys
from pathlib import Path

import numpy as np

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.config import SUCCESS_HORIZON  # noqa: E402
from ml.train_soy import (  # noqa: E402
    QUANTILES,
    QUANTILE_ROUNDS,
    _require_lightgbm,
    _train_booster,
    load_rows,
    split_xy,
)


def _percentile(values: np.ndarray, q: float) -> float:
    return float(np.quantile(values, q))


def _prf(y_true_pos: np.ndarray, y_pred_pos: np.ndarray) -> dict[str, float | int]:
    tp = int(np.sum(y_true_pos & y_pred_pos))
    fp = int(np.sum(~y_true_pos & y_pred_pos))
    fn = int(np.sum(y_true_pos & ~y_pred_pos))
    precision = tp / (tp + fp) if (tp + fp) else 0.0
    recall = tp / (tp + fn) if (tp + fn) else 0.0
    f1 = (2 * precision * recall / (precision + recall)) if (precision + recall) else 0.0
    return {
        "tp": tp,
        "fp": fp,
        "fn": fn,
        "precision": precision,
        "recall": recall,
        "f1": f1,
        "n_positive": int(np.sum(y_true_pos)),
        "n_alerts": int(np.sum(y_pred_pos)),
    }


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--export-dir",
        default=os.environ.get("ML_EXPORT_DIR", str(ROOT / ".local" / "ml" / "export")),
    )
    parser.add_argument(
        "--out-dir",
        default=os.environ.get("ML_REPORTS_DIR", str(ROOT / ".local" / "ml" / "reports")),
    )
    args = parser.parse_args()

    lgb = _require_lightgbm()
    parquet_path = Path(args.export_dir) / "soy_daily_training.parquet"
    if not parquet_path.is_file():
        print(f"missing export parquet: {parquet_path}", file=sys.stderr)
        return 1

    rows = load_rows(parquet_path)
    target = f"y_ret_{SUCCESS_HORIZON}d"
    x_train, y_train = split_xy(rows, target, "train")
    x_val, y_val = split_xy(rows, target, "val")
    x_test, y_test = split_xy(rows, target, "test")
    if len(y_train) < 50 or len(y_test) < 10:
        print("insufficient rows for early-warning eval", file=sys.stderr)
        return 1

    # Upside spike threshold: 90th percentile of train forward returns.
    spike_thr = _percentile(y_train, 0.90)
    abs_thr = _percentile(np.abs(y_train), 0.90)

    q_models: dict[str, object] = {}
    for q in QUANTILES:
        q_models[f"p{int(q * 100):02d}"] = _train_booster(
            lgb, x_train, y_train, x_val, y_val,
            objective="quantile", alpha=q, num_rounds=QUANTILE_ROUNDS,
        )

    p50 = np.asarray(q_models["p50"].predict(x_test))  # type: ignore[attr-defined]
    p90 = np.asarray(q_models["p90"].predict(x_test))  # type: ignore[attr-defined]
    p10 = np.asarray(q_models["p10"].predict(x_test))  # type: ignore[attr-defined]

    true_up = y_test >= spike_thr
    true_abs = np.abs(y_test) >= abs_thr

    rules = {
        "upside_p90_ge_train_p90": _prf(true_up, p90 >= spike_thr),
        "upside_p50_ge_train_p90": _prf(true_up, p50 >= spike_thr),
        "abs_extreme_p10_or_p90": _prf(true_abs, (p90 >= abs_thr) | (p10 <= -abs_thr)),
    }

    payload = {
        "target": target,
        "spike_threshold_train_p90": spike_thr,
        "abs_threshold_train_p90": abs_thr,
        "n_test": int(len(y_test)),
        "rules": rules,
        "rule_definitions": {
            "upside_p90_ge_train_p90": (
                "Alert when predicted p90 forward return ≥ train empirical p90; "
                "label = realized y_ret ≥ train p90"
            ),
            "upside_p50_ge_train_p90": (
                "Alert when predicted p50 ≥ train p90 (stricter median rule)"
            ),
            "abs_extreme_p10_or_p90": (
                "Alert when predicted p90 ≥ +|return|p90 or p10 ≤ −|return|p90; "
                "label = |realized| ≥ train |return| p90"
            ),
        },
    }

    out_dir = Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)
    json_path = out_dir / "soy_early_warning_eval.json"
    json_path.write_text(json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8")

    md = [
        "# Soy early-warning evaluation (Phase 31 C3)",
        "",
        f"Target: `{target}` · test n={len(y_test)}",
        "",
        f"- Upside spike threshold (train p90 of `{target}`): **{spike_thr:.6g}**",
        f"- Absolute-move threshold (train p90 of |return|): **{abs_thr:.6g}**",
        "",
        "## Rules",
        "",
        "| rule | precision | recall | f1 | tp | fp | fn | alerts |",
        "|------|-----------|--------|----|----|----|----|--------|",
    ]
    for name, stats in rules.items():
        md.append(
            f"| {name} | {stats['precision']:.3f} | {stats['recall']:.3f} | {stats['f1']:.3f} | "
            f"{stats['tp']} | {stats['fp']} | {stats['fn']} | {stats['n_alerts']} |"
        )
    md.extend(["", "## Definitions", ""])
    for name, text in payload["rule_definitions"].items():  # type: ignore[union-attr]
        md.append(f"- **{name}:** {text}")
    md_path = out_dir / "soy_early_warning_eval.md"
    md_path.write_text("\n".join(md) + "\n", encoding="utf-8")
    print(f"wrote {json_path}")
    print(f"wrote {md_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
