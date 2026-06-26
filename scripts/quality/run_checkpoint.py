#!/usr/bin/env python3
"""Run a Great Expectations bronze checkpoint against local lake Parquet."""

from __future__ import annotations

import argparse
import json
import sys
import uuid
from pathlib import Path

import great_expectations as gx
import pandas as pd
from great_expectations.core import ExpectationSuite
from great_expectations.exceptions import DataContextError

REPO_ROOT = Path(__file__).resolve().parents[2]
DEFAULT_CONTEXT_ROOT = REPO_ROOT / "expectations" / "gx"


def load_bronze_dataframe(bronze_dir: Path) -> pd.DataFrame:
    files = sorted(bronze_dir.glob("ingest_date=*/part-*.parquet"))
    if not files:
        raise FileNotFoundError(f"no bronze parquet under {bronze_dir}")
    frames = [pd.read_parquet(path) for path in files]
    return pd.concat(frames, ignore_index=True)


def sync_suite_to_context(context_root: Path, suite_path: Path) -> str:
    payload = json.loads(suite_path.read_text(encoding="utf-8"))
    suite_name = payload["name"]
    suite = ExpectationSuite(**payload)

    try:
        context = gx.get_context(project_root_dir=str(context_root))
        try:
            context.suites.delete(suite_name)
        except (DataContextError, KeyError, ValueError):
            pass
        context.suites.add(suite)
        return suite_name
    except Exception:
        # Fallback: validate without persisting to context store.
        return suite_name


def validate_bronze(
    context_root: Path,
    suite_path: Path,
    bronze_dir: Path,
) -> dict:
    df = load_bronze_dataframe(bronze_dir)
    payload = json.loads(suite_path.read_text(encoding="utf-8"))
    suite = ExpectationSuite(**payload)

    context = gx.get_context(project_root_dir=str(context_root))
    datasource_name = f"bronze_runtime_{bronze_dir.as_posix().replace('/', '_')[-48:]}"
    try:
        context.data_sources.get(datasource_name)
    except (DataContextError, KeyError, ValueError):
        context.data_sources.add_pandas(datasource_name)

    batch = context.data_sources.get(datasource_name).read_dataframe(
        df,
        asset_name=f"bronze_batch_{uuid.uuid4().hex[:8]}",
    )
    validator = context.get_validator(batch=batch, expectation_suite=suite)
    result = validator.validate()
    return {
        "success": result.success,
        "evaluated_expectations": result.statistics.get("evaluated_expectations", 0),
        "successful_expectations": result.statistics.get("successful_expectations", 0),
        "unsuccessful_expectations": result.statistics.get("unsuccessful_expectations", 0),
        "bronze_dir": str(bronze_dir),
        "suite": payload["name"],
    }


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate bronze Parquet with Great Expectations")
    parser.add_argument(
        "--checkpoint",
        default="bronze_conab_estimativa_graos",
        help="Checkpoint name (maps to suite file when --suite omitted)",
    )
    parser.add_argument("--suite", help="Path to expectation suite JSON")
    parser.add_argument("--bronze-dir", required=True, help="Bronze dataset directory")
    parser.add_argument(
        "--context-root",
        default=str(DEFAULT_CONTEXT_ROOT),
        help="Great Expectations project root",
    )
    args = parser.parse_args()

    context_root = Path(args.context_root).resolve()
    bronze_dir = Path(args.bronze_dir).resolve()

    if args.suite:
        suite_path = Path(args.suite).resolve()
    else:
        mapping = {
            "bronze_conab_estimativa_graos": "expectations/suites/bronze/conab/estimativa_graos.json",
            "bronze_conab_estimativa_cana": "expectations/suites/bronze/conab/estimativa_cana.json",
            "bronze_conab_serie_historica_graos": "expectations/suites/bronze/conab/serie_historica_graos.json",
            "bronze_conab_serie_historica_cana": "expectations/suites/bronze/conab/serie_historica_cana.json",
            "bronze_conab_oferta_demanda": "expectations/suites/bronze/conab/oferta_demanda.json",
            "bronze_conab_precos_minimos": "expectations/suites/bronze/conab/precos_minimos.json",
            "bronze_conab_precos_semanal_uf": "expectations/suites/bronze/conab/precos_semanal_uf.json",
            "bronze_conab_precos_semanal_municipio": "expectations/suites/bronze/conab/precos_semanal_municipio.json",
            "bronze_conab_precos_mensal_uf": "expectations/suites/bronze/conab/precos_mensal_uf.json",
            "bronze_conab_precos_mensal_municipio": "expectations/suites/bronze/conab/precos_mensal_municipio.json",
            "bronze_conab_prohort_diario": "expectations/suites/bronze/conab/prohort_diario.json",
            "bronze_conab_prohort_mensal": "expectations/suites/bronze/conab/prohort_mensal.json",
            "bronze_conab_estoques_publicos": "expectations/suites/bronze/conab/estoques_publicos.json",
            "bronze_conab_operacoes_comercializacao": "expectations/suites/bronze/conab/operacoes_comercializacao.json",
            "bronze_conab_vendas_balcao": "expectations/suites/bronze/conab/vendas_balcao.json",
            "bronze_anp_combustiveis_precos_medios_municipios": "expectations/suites/bronze/anp/combustiveis_precos_medios_municipios.json",
            "bronze_anp_combustiveis_precos_postos": "expectations/suites/bronze/anp/combustiveis_precos_postos.json",
            "bronze_conab_armazenagem": "expectations/suites/bronze/conab/armazenagem.json",
            "bronze_conab_frete": "expectations/suites/bronze/conab/frete.json",
            "bronze_conab_serie_historica_capacidade_estatica": "expectations/suites/bronze/conab/serie_historica_capacidade_estatica.json",
            "bronze_conab_alimenta_brasil_entregas": "expectations/suites/bronze/conab/alimenta_brasil_entregas.json",
            "bronze_conab_alimenta_brasil_propostas": "expectations/suites/bronze/conab/alimenta_brasil_propostas.json",
            "bronze_ibge_localidades_municipios": "expectations/suites/bronze/ibge/localidades_municipios.json",
            "bronze_ibge_localidades_ufs": "expectations/suites/bronze/ibge/localidades_ufs.json",
            "bronze_ibge_localidades_regioes": "expectations/suites/bronze/ibge/localidades_regioes.json",
            "bronze_ibge_localidades_mesorregioes": "expectations/suites/bronze/ibge/localidades_mesorregioes.json",
            "bronze_ibge_localidades_microrregioes": "expectations/suites/bronze/ibge/localidades_microrregioes.json",
            "bronze_ibge_pam_area_quantidade": "expectations/suites/bronze/ibge/pam_area_quantidade.json",
            "bronze_ibge_pam_rendimento_valor": "expectations/suites/bronze/ibge/pam_rendimento_valor.json",
            "bronze_ibge_pam_estabelecimentos": "expectations/suites/bronze/ibge/pam_estabelecimentos.json",
            "bronze_inmet_estacoes_automaticas": "expectations/suites/bronze/inmet/estacoes_automaticas.json",
            "bronze_inmet_estacoes_convencionais": "expectations/suites/bronze/inmet/estacoes_convencionais.json",
            "bronze_inmet_bdmep_diario": "expectations/suites/bronze/inmet/bdmep_diario.json",
            "bronze_inmet_bdmep_mensal": "expectations/suites/bronze/inmet/bdmep_mensal.json",
            "bronze_bcb_sgs_ipca": "expectations/suites/bronze/bcb/sgs_ipca.json",
            "bronze_bcb_sgs_ptax_usd_venda": "expectations/suites/bronze/bcb/sgs_ptax_usd_venda.json",
            "bronze_bcb_sgs_ipca_12m": "expectations/suites/bronze/bcb/sgs_ipca_12m.json",
            "bronze_bcb_sgs_igpm": "expectations/suites/bronze/bcb/sgs_igpm.json",
            "bronze_bcb_sgs_ptax_usd_compra": "expectations/suites/bronze/bcb/sgs_ptax_usd_compra.json",
            "bronze_cepea_soja_paranagua": "expectations/suites/bronze/cepea/soja_paranagua.json",
            "bronze_cepea_soja_parana": "expectations/suites/bronze/cepea/soja_parana.json",
            "bronze_cepea_milho": "expectations/suites/bronze/cepea/milho.json",
        }
        rel = mapping.get(args.checkpoint)
        if not rel:
            print(f"unknown checkpoint {args.checkpoint}", file=sys.stderr)
            return 2
        suite_path = (REPO_ROOT / rel).resolve()

    if not suite_path.is_file():
        print(f"suite not found: {suite_path}", file=sys.stderr)
        return 2

    try:
        sync_suite_to_context(context_root, suite_path)
        summary = validate_bronze(context_root, suite_path, bronze_dir)
    except Exception as exc:  # noqa: BLE001 — CLI boundary
        print(json.dumps({"success": False, "error": str(exc)}))
        return 1

    print(json.dumps(summary))
    return 0 if summary["success"] else 1


if __name__ == "__main__":
    sys.exit(main())
