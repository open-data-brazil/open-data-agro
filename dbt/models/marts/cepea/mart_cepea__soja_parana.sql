{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_cepea__soja_parana/mart.parquet',
    format='parquet'
) }}

select
    produto,
    praca,
    data,
    preco_rs_sc,
    variacao_dia_pct,
    preco_usd_sc,
    ano,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_cepea__soja_parana') }}
