{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ipea__series_macro_regionais/mart.parquet',
    format='parquet'
) }}

select
    series_code,
    refdate,
    value,
    region_level,
    territory_code,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_ipea__series_macro_regionais') }}
