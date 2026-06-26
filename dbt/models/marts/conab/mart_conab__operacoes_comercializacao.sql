{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__operacoes_comercializacao/mart.parquet',
    format='parquet'
) }}

select
    id_edital,
    num_lote,
    num_dco,
    dsc_dco,
    situacao_dco,
    produto,
    id_produto,
    tipo_operacao,
    operacao,
    ano_edital,
    mes_edital,
    uf_armazem_origem,
    unidade_comercializacao,
    qtd_ofertada,
    qtd_negociada,
    vlr_operacao_s_icms,
    unidade_medida_ofertada_negociada,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__operacoes_comercializacao') }}
