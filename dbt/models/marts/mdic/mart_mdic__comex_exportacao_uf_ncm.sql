{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_mdic__comex_exportacao_uf_ncm/mart.parquet',
    format='parquet'
) }}

select
    co_ncm,
    ncm_descricao,
    produto_slug,
    uf,
    data,
    valor_fob_usd,
    quantidade_kg,
    ano,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_mdic__comex_exportacao_uf_ncm') }}
