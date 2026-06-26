{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_argentina__magyp_producion_granos/mart.parquet',
    format='parquet'
) }}

select
    series_id,
    commodity_slug,
    refyear,
    value,
    unit,
    source,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_argentina__magyp_producion_granos') }}
