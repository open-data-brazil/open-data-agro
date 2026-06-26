{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_eurostat__ag_prices/mart.parquet',
    format='parquet'
) }}

select
    dataset_code,
    geo,
    product_code,
    product_name,
    year,
    index_value,
    base_period,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_eurostat__ag_prices') }}
