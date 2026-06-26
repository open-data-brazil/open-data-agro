{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_bcb__sgs_ptax_usd_venda/mart.parquet',
    format='parquet'
) }}

select
    sgs_codigo,
    data,
    valor,
    ano,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_bcb__sgs_ptax_usd_venda') }}
