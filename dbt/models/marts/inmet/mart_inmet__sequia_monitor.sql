{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_inmet__sequia_monitor/mart.parquet',
    format='parquet'
) }}

select
    mapa_id,
    ano,
    mes,
    categoria_seca,
    area_km2,
    area_id,
    tipo_area,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_inmet__sequia_monitor') }}
