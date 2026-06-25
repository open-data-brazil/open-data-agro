{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__oferta_demanda/mart.parquet',
    format='parquet'
) }}

select
    produto,
    safra,
    id_produto,
    estoque_inicial_1000t,
    producao_1000t,
    importacao_1000t,
    consumo_1000t,
    exportacao_1000t,
    estoque_final_1000t,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__oferta_demanda') }}
