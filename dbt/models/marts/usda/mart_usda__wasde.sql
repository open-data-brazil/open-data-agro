{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_usda__wasde/mart.parquet',
    format='parquet'
) }}

select
    report_month,
    commodity,
    market_year,
    attribute,
    value,
    unit,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_usda__wasde') }}
