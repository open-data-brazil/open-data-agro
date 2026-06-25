{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__alimenta_brasil_entregas/mart.parquet',
    format='parquet'
) }}

select
    ano_entrega,
    mes_entrega,
    municipio,
    uf,
    regiao,
    sexo,
    unidade_medida,
    qtd_entregue,
    valor_entregue,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__alimenta_brasil_entregas') }}
