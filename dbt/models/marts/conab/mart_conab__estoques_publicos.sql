{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__estoques_publicos/mart.parquet',
    format='parquet'
) }}

select
    produto,
    id_produto,
    municipio,
    cod_ibge,
    uf,
    ano,
    mes,
    conta_operacional,
    qtd_estoque_kg,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__estoques_publicos') }}
