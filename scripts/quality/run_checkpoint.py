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
            "bronze_conab_estimativa_cafe": "expectations/suites/bronze/conab/estimativa_cafe.json",
            "bronze_conab_serie_historica_cafe": "expectations/suites/bronze/conab/serie_historica_cafe.json",
            "bronze_conab_custo_producao": "expectations/suites/bronze/conab/custo_producao.json",
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
            "bronze_anp_etanol_precos": "expectations/suites/bronze/anp/etanol_precos.json",
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
            "bronze_ibge_lspa_area_producao": "expectations/suites/bronze/ibge/lspa_area_producao.json",
            "bronze_inmet_estacoes_automaticas": "expectations/suites/bronze/inmet/estacoes_automaticas.json",
            "bronze_inmet_estacoes_convencionais": "expectations/suites/bronze/inmet/estacoes_convencionais.json",
            "bronze_inmet_bdmep_diario": "expectations/suites/bronze/inmet/bdmep_diario.json",
            "bronze_inmet_bdmep_mensal": "expectations/suites/bronze/inmet/bdmep_mensal.json",
            "bronze_inmet_pacote_anual_automaticas": "expectations/suites/bronze/inmet/pacote_anual_automaticas.json",
            "bronze_bcb_sgs_ipca": "expectations/suites/bronze/bcb/sgs_ipca.json",
            "bronze_bcb_sgs_ptax_usd_venda": "expectations/suites/bronze/bcb/sgs_ptax_usd_venda.json",
            "bronze_bcb_sgs_ipca_12m": "expectations/suites/bronze/bcb/sgs_ipca_12m.json",
            "bronze_bcb_sgs_igpm": "expectations/suites/bronze/bcb/sgs_igpm.json",
            "bronze_bcb_sgs_ptax_usd_compra": "expectations/suites/bronze/bcb/sgs_ptax_usd_compra.json",
            "bronze_bcb_sgs_selic": "expectations/suites/bronze/bcb/sgs_selic.json",
            "bronze_bcb_cim_agro_credito_rural": "expectations/suites/bronze/bcb/cim_agro_credito_rural.json",
            "bronze_cepea_soja_paranagua": "expectations/suites/bronze/cepea/soja_paranagua.json",
            "bronze_cepea_soja_parana": "expectations/suites/bronze/cepea/soja_parana.json",
            "bronze_cepea_milho": "expectations/suites/bronze/cepea/milho.json",
            "bronze_cepea_boi_gordo": "expectations/suites/bronze/cepea/boi_gordo.json",
            "bronze_mdic_comex_exportacao_ncm_mes": "expectations/suites/bronze/mdic/comex_exportacao_ncm_mes.json",
            "bronze_mdic_comex_importacao_ncm_mes": "expectations/suites/bronze/mdic/comex_importacao_ncm_mes.json",
            "bronze_mdic_comex_exportacao_uf_ncm": "expectations/suites/bronze/mdic/comex_exportacao_uf_ncm.json",
            "bronze_mdic_comex_importacao_diesel_ncm": "expectations/suites/bronze/mdic/comex_importacao_diesel_ncm.json",
            "bronze_antt_pracas_pedagio": "expectations/suites/bronze/antt/pracas_pedagio.json",
            "bronze_mapa_zarc_tabua_risco": "expectations/suites/bronze/mapa/zarc_tabua_risco.json",
            "bronze_mapa_agrofit_produtos_formulados": "expectations/suites/bronze/mapa/agrofit_produtos_formulados.json",
            "bronze_mapa_agrofit_produtos_tecnicos": "expectations/suites/bronze/mapa/agrofit_produtos_tecnicos.json",
            "bronze_mapa_sipeagro_estabelecimentos": "expectations/suites/bronze/mapa/sipeagro_estabelecimentos.json",
            "bronze_mapa_sipeagro_produtos": "expectations/suites/bronze/mapa/sipeagro_produtos.json",
            "bronze_mapa_sigef_producao_sementes": "expectations/suites/bronze/mapa/sigef_producao_sementes.json",
            "bronze_mapa_sigef_areas": "expectations/suites/bronze/mapa/sigef_areas.json",
            "bronze_mapa_sisser_seguro_rural": "expectations/suites/bronze/mapa/sisser_seguro_rural.json",
            "bronze_ibama_sisfogo_incendios": "expectations/suites/bronze/ibama/sisfogo_incendios.json",
            "bronze_ibama_licencas_ambientais": "expectations/suites/bronze/ibama/licencas_ambientais.json",
            "bronze_ibama_autos_infracao": "expectations/suites/bronze/ibama/autos_infracao.json",
            "bronze_ana_pluviometria_redes": "expectations/suites/bronze/ana/pluviometria_redes.json",
            "bronze_embrapa_agroapi_agrofit": "expectations/suites/bronze/embrapa/agroapi_agrofit.json",
            "bronze_transportes_mtr_bit_malha_shapefile": "expectations/suites/bronze/transportes/mtr_bit_malha_shapefile.json",
            "bronze_ana_hidrologia_series": "expectations/suites/bronze/ana/hidrologia_series.json",
            "bronze_dnit_snv_rodovias_federais": "expectations/suites/bronze/dnit/snv_rodovias_federais.json",
            "bronze_ipea_series_macro_regionais": "expectations/suites/bronze/ipea/series_macro_regionais.json",
            "bronze_ibge_pevs_producao_vegetal": "expectations/suites/bronze/ibge/pevs_producao_vegetal.json",
            "bronze_ibge_ppm_producao_municipal": "expectations/suites/bronze/ibge/ppm_producao_municipal.json",
            "bronze_ibge_ppm_efetivo_rebanhos": "expectations/suites/bronze/ibge/ppm_efetivo_rebanhos.json",
            "bronze_ibge_ppm_vacas_ordenhadas": "expectations/suites/bronze/ibge/ppm_vacas_ordenhadas.json",
            "bronze_ibge_ppm_ovinos_tosquiados": "expectations/suites/bronze/ibge/ppm_ovinos_tosquiados.json",
            "bronze_ibge_ppm_aquicultura": "expectations/suites/bronze/ibge/ppm_aquicultura.json",
            "bronze_ibge_pam_precos_produtor": "expectations/suites/bronze/ibge/pam_precos_produtor.json",
            "bronze_ibge_pam_culturas_estendidas": "expectations/suites/bronze/ibge/pam_culturas_estendidas.json",
            "bronze_ibge_lspa_rendimento_medio": "expectations/suites/bronze/ibge/lspa_rendimento_medio.json",
            "bronze_ibge_censo_agro_area_uso_solo": "expectations/suites/bronze/ibge/censo_agro_area_uso_solo.json",
            "bronze_ibge_censo_agro_maquinario": "expectations/suites/bronze/ibge/censo_agro_maquinario.json",
            "bronze_ibge_pnad_rural_renda_ocupacao": "expectations/suites/bronze/ibge/pnad_rural_renda_ocupacao.json",
            "bronze_aneel_tarifas_energia": "expectations/suites/bronze/aneel/tarifas_energia.json",
            "bronze_bndes_financiamento_agro": "expectations/suites/bronze/bndes/financiamento_agro.json",
            "bronze_bndes_desembolsos_linhas_agro": "expectations/suites/bronze/bndes/desembolsos_linhas_agro.json",
            "bronze_inmet_sequia_monitor": "expectations/suites/bronze/inmet/sequia_monitor.json",
            "bronze_oecd_ag_outlook": "expectations/suites/bronze/oecd/ag_outlook.json",
            "bronze_fao_food_price_index": "expectations/suites/bronze/fao/food_price_index.json",
            "bronze_argentina_magyp_producion_granos": "expectations/suites/bronze/argentina/magyp_producion_granos.json",
            "bronze_b3_futuro_soja": "expectations/suites/bronze/b3/futuro_soja.json",
            "bronze_b3_futuro_milho": "expectations/suites/bronze/b3/futuro_milho.json",
            "bronze_b3_futuro_boi": "expectations/suites/bronze/b3/futuro_boi.json",
            "bronze_fao_prices_agro": "expectations/suites/bronze/fao/prices_agro.json",
            "bronze_fao_producao_agro": "expectations/suites/bronze/fao/producao_agro.json",
            "bronze_worldbank_pink_sheet_monthly": "expectations/suites/bronze/worldbank/pink_sheet_monthly.json",
            "bronze_worldbank_ag_indices": "expectations/suites/bronze/worldbank/ag_indices.json",
            "bronze_noaa_enso_indices": "expectations/suites/bronze/noaa/enso_indices.json",
            "bronze_noaa_global_temp_anomaly": "expectations/suites/bronze/noaa/global_temp_anomaly.json",
            "bronze_eia_petroleum_prices": "expectations/suites/bronze/eia/petroleum_prices.json",
            "bronze_usda_wasde": "expectations/suites/bronze/usda/wasde.json",
            "bronze_igc_goi_index": "expectations/suites/bronze/igc/goi_index.json",
            "bronze_eurostat_ag_prices": "expectations/suites/bronze/eurostat/ag_prices.json",
            "bronze_argentina_bcra_cambio": "expectations/suites/bronze/argentina/bcra_cambio.json",
            "bronze_un_comtrade_bulk": "expectations/suites/bronze/un/comtrade_bulk.json",
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
