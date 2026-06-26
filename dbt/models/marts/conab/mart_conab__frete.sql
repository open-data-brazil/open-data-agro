{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__frete/mart.parquet',
    format='parquet'
) }}

select
    fonte,
    municipio_origem,
    cod_ibge_origem,
    uf_origem,
    municipio_destino,
    cod_ibge_destino,
    uf_destino,
    ano,
    mes,
    distancia_km,
    valor_frete_tonelada,
    valor_tonelada_km,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__frete') }}
