{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_antaq__movimentacao_carga_portuaria/mart.parquet',
    format='parquet'
) }}

select
    ano,
    mes,
    codigo_instalacao_portuaria,
    nome_instalacao_portuaria,
    tipo_movimentacao,
    tipo_navegacao,
    sentido,
    natureza_carga,
    tipo_operacao,
    peso_toneladas,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_antaq__movimentacao_carga_portuaria') }}
