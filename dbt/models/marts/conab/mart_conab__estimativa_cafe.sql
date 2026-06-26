{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__estimativa_cafe/mart.parquet',
    format='parquet'
) }}

select
    produto,
    uf,
    safra,
    tipo_safra,
    id_produto,
    id_levantamento,
    levantamento,
    area_plantada_mil_ha,
    producao_mil_t,
    produtividade_mil_ha_mil_t,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__estimativa_cafe') }}
