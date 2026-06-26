{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_mdic__comex_importacao_ncm_mes/mart.parquet',
    format='parquet'
) }}

select
    co_ncm,
    ncm_descricao,
    produto_slug,
    data,
    valor_cif_usd,
    quantidade_kg,
    valor_frete_usd,
    valor_seguro_usd,
    ano,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_mdic__comex_importacao_ncm_mes') }}
