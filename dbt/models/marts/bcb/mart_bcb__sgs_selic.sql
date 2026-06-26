{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_bcb__sgs_selic/mart.parquet',
    format='parquet'
) }}

select
    sgs_codigo,
    data,
    valor,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_bcb__sgs_selic') }}
