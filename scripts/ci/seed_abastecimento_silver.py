#!/usr/bin/env python3
"""Seed minimal silver Delta for CONAB Abastecimento + ANP combustíveis CI (local-first)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import pyarrow as pa
from deltalake import write_deltalake


def write_table(root: Path, agency: str, table: str, data: pa.Table) -> None:
    path = root / "silver" / agency / table
    path.parent.mkdir(parents=True, exist_ok=True)
    write_deltalake(str(path), data, mode="overwrite")


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    for mart in (
        "mart_conab__estoques_publicos",
        "mart_conab__operacoes_comercializacao",
        "mart_conab__vendas_balcao",
        "mart_anp__combustiveis_precos_medios_municipios",
        "mart_anp__combustiveis_precos_postos",
    ):
        (lake_root / "gold" / mart).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-25T12:00:00Z"

    estoques = pa.table(
        {
            "produto": ["ARROZ", "MILHO"],
            "id_produto": ["4693", "4694"],
            "nom_municipio": ["ANANINDEUA-PA", "CUIABA-MT"],
            "cod_ibge": ["1500800", "5103403"],
            "uf": ["PA", "MT"],
            "num_ano": ["2022", "2022"],
            "num_mes": ["1", "1"],
            "conta_operacional": ["ESTRATEGIC", "ESTRATEGIC"],
            "qtd_estoque_kg": ["660,0", "1200,0"],
            "_dataset_id": ["conab.estoques-publicos", "conab.estoques-publicos"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "conab", "estoques_publicos", estoques)

    operacoes = pa.table(
        {
            "id_edital": ["SEC-201313", "SEC-201315"],
            "num_lote": ["1", "1"],
            "num_dco": ["0047252936", "0047252944"],
            "dsc_dco": ["Autorização de Venda Terceiros", "Autorização de Venda Terceiros"],
            "dsc_situacao_dco": ["NÃO PAGO", "NÃO PAGO"],
            "produto": ["SOJA", "SOJA"],
            "id_produto": ["4744", "4744"],
            "dsc_tipo_operacao": ["VENDA", "VENDA"],
            "dsc_operacao": ["PRIVADA", "PRIVADA"],
            "ano_edital": ["2013", "2013"],
            "mes_edital": ["8", "8"],
            "uf_armazem_origem": ["MT", "MT"],
            "dsc_unidade_comercializacao": ["1 KG (GRANEL)", "1 KG (GRANEL)"],
            "qtd_ofertada": ["60000.0", "17000.0"],
            "qtd_negociada": ["60000.0", "17000.0"],
            "vlr_operacao_s_icms": ["55800.0", "15810.0"],
            "dsc_unidade_medida_ofertada_negociada": ["QUILO", "QUILO"],
            "_dataset_id": [
                "conab.operacoes-comercializacao",
                "conab.operacoes-comercializacao",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "conab", "operacoes_comercializacao", operacoes)

    vendas = pa.table(
        {
            "num_ano_gravacao": ["2025", "2026"],
            "num_mes_gravacao": ["11", "1"],
            "munipio_armazem_venda": ["RONDONÓPOLIS", "RONDONÓPOLIS"],
            "uf": ["MT", "MT"],
            "qtd_produto_kg": ["7820", "9560"],
            "valor_comercializado": ["7141,112", "10711,98"],
            "numero_atendimentos": ["6", "8"],
            "clientes_atendidos": ["6", "8"],
            "_dataset_id": ["conab.vendas-balcao", "conab.vendas-balcao"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "conab", "vendas_balcao", vendas)

    medios = pa.table(
        {
            "DATA INICIAL": ["2026-06-14 00:00:00", "2026-06-14 00:00:00"],
            "DATA FINAL": ["2026-06-20 00:00:00", "2026-06-20 00:00:00"],
            "ESTADO": ["MATO GROSSO DO SUL", "SAO PAULO"],
            "MUNICÍPIO": ["CAMPO GRANDE", "SAO PAULO"],
            "PRODUTO": ["GASOLINA COMUM", "ETANOL HIDRATADO"],
            "NÚMERO DE POSTOS PESQUISADOS": ["45", "120"],
            "UNIDADE DE MEDIDA": ["R$/l", "R$/l"],
            "PREÇO MÉDIO REVENDA": ["5.89", "4.12"],
            "DESVIO PADRÃO REVENDA": ["0.21", "0.55"],
            "PREÇO MÍNIMO REVENDA": ["5.49", "3.89"],
            "PREÇO MÁXIMO REVENDA": ["6.29", "5.10"],
            "COEF DE VARIAÇÃO REVENDA": ["0.036", "0.133"],
            "_dataset_id": [
                "anp.combustiveis-precos-medios-municipios",
                "anp.combustiveis-precos-medios-municipios",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "anp", "combustiveis_precos_medios_municipios", medios)

    postos = pa.table(
        {
            "CNPJ": ["61602199002409", "00000000000191"],
            "RAZÃO": ["COMPANHIA ULTRAGAZ S A", "POSTO EXEMPLO LTDA"],
            "FANTASIA": ["ULTRAGAZ", "POSTO EXEMPLO"],
            "ENDEREÇO": ["RUA AMARO CASTRO LIMA", "AV PAULISTA"],
            "NÚMERO": ["1852", "1000"],
            "COMPLEMENTO": ["", ""],
            "BAIRRO": ["VILA NOVA CAMPO GRANDE", "BELA VISTA"],
            "CEP": ["79106361", "01310100"],
            "MUNICÍPIO": ["CAMPO GRANDE", "SAO PAULO"],
            "ESTADO": ["MATO GROSSO DO SUL", "SAO PAULO"],
            "BANDEIRA": ["ULTRAGAZ", "BR"],
            "PRODUTO": ["GLP", "GASOLINA COMUM"],
            "UNIDADE DE MEDIDA": ["R$/kg", "R$/l"],
            "PREÇO DE REVENDA": ["99.90", "5.79"],
            "DATA DA COLETA": ["2026-06-18", "2026-06-19"],
            "_dataset_id": ["anp.combustiveis-precos-postos", "anp.combustiveis-precos-postos"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "anp", "combustiveis_precos_postos", postos)

    print(f"seeded abastecimento + ANP silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
