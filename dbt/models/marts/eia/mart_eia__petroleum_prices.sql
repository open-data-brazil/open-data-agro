{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_eia__petroleum_prices/mart.parquet',
    format='parquet'
) }}

select
    series_id,
    series_name,
    commodity_slug,
    refdate,
    unit,
    value,
    frequency,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_eia__petroleum_prices') }}
