#!/usr/bin/env python3
"""SHAP + ablation for frete LightGBM (Phase 32 F4)."""

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

from ml.config_frete import NUMERIC_FEATURE_COLS  # noqa: E402
from ml.metrics import mae  # noqa: E402
from ml.train_frete import FEATURE_COLS, _require_lgb, _train, split_xy  # noqa: E402

ABLATION_GROUPS = {
    "distance": ("distancia_km",),
    "calendar": ("mes", "ano"),
    "corridor_lags": ("corridor_lag12_ton_avg", "corridor_lag1_ton_avg"),
    "route_ids": ("uf_origem", "uf_destino", "fonte_code"),
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
    if not path.is_file():
        print(f"missing {path}", file=sys.stderr)
        return 1

    rows = pq.read_table(path).to_pylist()
    x_tr, y_tr = split_xy(rows, "train")
    x_va, y_va = split_xy(rows, "val")
    x_te, y_te = split_xy(rows, "test")
    model = _train(lgb, x_tr, y_tr, x_va, y_va, objective="regression", alpha=None, rounds=600)
    full_mae = mae(y_te.tolist(), model.predict(x_te).tolist())
    assert full_mae is not None

    import shap

    n = min(200, len(x_te))
    rng = np.random.default_rng(42)
    idx = rng.choice(len(x_te), size=n, replace=False)
    values = np.asarray(shap.TreeExplainer(model).shap_values(x_te[idx]))
    means = np.mean(np.abs(values), axis=0)
    shap_ranked = sorted(
        (
            {"feature": FEATURE_COLS[i], "mean_abs_shap": float(means[i])}
            for i in range(len(FEATURE_COLS))
        ),
        key=lambda d: float(d["mean_abs_shap"]),
        reverse=True,
    )

    # Ablation: zero-out feature groups on test preds via retrain without those cols is heavy;
    # use permutation importance on test instead for speed.
    ablation: dict[str, object] = {}
    base_pred = model.predict(x_te)
    for name, cols in ABLATION_GROUPS.items():
        x_perm = x_te.copy()
        for col in cols:
            if col not in FEATURE_COLS:
                continue
            j = FEATURE_COLS.index(col)
            x_perm[:, j] = rng.permutation(x_perm[:, j])
        perm_mae = mae(y_te.tolist(), model.predict(x_perm).tolist())
        assert perm_mae is not None
        ablation[name] = {
            "permuted_features": list(cols),
            "test_mae": perm_mae,
            "mae_delta_vs_full": perm_mae - full_mae,
        }

    payload = {
        "target": "valor_frete_tonelada",
        "full_model_test_mae": full_mae,
        "shap_sample_size": n,
        "shap_mean_abs": shap_ranked,
        "ablation_permutation": ablation,
        "numeric_features": list(NUMERIC_FEATURE_COLS),
    }
    out = Path(args.out_dir)
    out.mkdir(parents=True, exist_ok=True)
    (out / "frete_shap_ablation.json").write_text(
        json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
    md = [
        "# Frete SHAP + ablation (Phase 32 F4)",
        "",
        f"Full-model test MAE: **{full_mae:.4f}** · SHAP n={n}",
        "",
        "| feature | mean_abs_shap |",
        "|---------|---------------|",
    ]
    for row in shap_ranked:
        md.append(f"| {row['feature']} | {row['mean_abs_shap']:.6g} |")
    md.extend(
        ["", "## Permutation ablation", "", "| group | test MAE | Δ |", "|-------|----------|---|"]
    )
    for name, stats in ablation.items():
        md.append(
            f"| {name} | {stats['test_mae']:.4f} | {stats['mae_delta_vs_full']:+.4f} |"  # type: ignore[index]
        )
    md_path = out / "frete_shap_ablation.md"
    md_path.write_text("\n".join(md) + "\n", encoding="utf-8")
    print(f"wrote {md_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
