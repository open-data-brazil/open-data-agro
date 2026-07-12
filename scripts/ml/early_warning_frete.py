#!/usr/bin/env python3
"""High-frete alert precision/recall on test (Phase 32 F4)."""

from __future__ import annotations

import argparse
import json
import os
import sys
from pathlib import Path

import numpy as np
import pyarrow.parquet as pq

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.train_frete import QUANTILES, _require_lgb, _train, split_xy  # noqa: E402


def _prf(y_true: np.ndarray, y_pred: np.ndarray) -> dict[str, float | int]:
    tp = int(np.sum(y_true & y_pred))
    fp = int(np.sum(~y_true & y_pred))
    fn = int(np.sum(y_true & ~y_pred))
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
        "n_positive": int(np.sum(y_true)),
        "n_alerts": int(np.sum(y_pred)),
    }


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--export-dir",
        default=os.environ.get(
            "ML_FRETE_EXPORT_DIR", str(ROOT / ".local" / "ml" / "export_frete")
        ),
    )
    parser.add_argument(
        "--out-dir",
        default=os.environ.get("ML_REPORTS_DIR", str(ROOT / ".local" / "ml" / "reports")),
    )
    args = parser.parse_args()

    lgb = _require_lgb()
    path = Path(args.export_dir) / "frete_training.parquet"
    rows = pq.read_table(path).to_pylist()
    x_tr, y_tr = split_xy(rows, "train")
    x_va, y_va = split_xy(rows, "val")
    x_te, y_te = split_xy(rows, "test")
    thr = float(np.quantile(y_tr, 0.90))

    models = {
        f"p{int(q * 100):02d}": _train(
            lgb, x_tr, y_tr, x_va, y_va, objective="quantile", alpha=q, rounds=400
        )
        for q in QUANTILES
    }
    p50 = np.asarray(models["p50"].predict(x_te))
    p90 = np.asarray(models["p90"].predict(x_te))
    true_hi = y_te >= thr
    rules = {
        "high_frete_p90_ge_train_p90": _prf(true_hi, p90 >= thr),
        "high_frete_p50_ge_train_p90": _prf(true_hi, p50 >= thr),
    }

    payload = {
        "target": "valor_frete_tonelada",
        "high_frete_threshold_train_p90": thr,
        "n_test": int(len(y_te)),
        "rules": rules,
    }
    out = Path(args.out_dir)
    out.mkdir(parents=True, exist_ok=True)
    (out / "frete_early_warning_eval.json").write_text(
        json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
    md = [
        "# Frete high-cost early-warning (Phase 32)",
        "",
        f"Threshold (train p90): **{thr:.2f}** R$/t · test n={len(y_te)}",
        "",
        "| rule | precision | recall | f1 | alerts |",
        "|------|-----------|--------|----|--------|",
    ]
    for name, s in rules.items():
        md.append(
            f"| {name} | {s['precision']:.3f} | {s['recall']:.3f} | {s['f1']:.3f} | {s['n_alerts']} |"
        )
    md_path = out / "frete_early_warning_eval.md"
    md_path.write_text("\n".join(md) + "\n", encoding="utf-8")
    print(f"wrote {md_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
