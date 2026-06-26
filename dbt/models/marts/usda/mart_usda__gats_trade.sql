{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_usda__gats_trade/mart.parquet',
    format='parquet'
) }}

select
    commodity_code,
    commodity_name,
    partner_code,
    partner_name,
    flow,
    year,
    value,
    unit,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_usda__gats_trade') }}
