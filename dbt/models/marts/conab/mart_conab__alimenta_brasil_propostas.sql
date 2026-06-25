{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__alimenta_brasil_propostas/mart.parquet',
    format='parquet'
) }}

select
    ano,
    mes,
    municipio,
    cod_ibge,
    uf,
    regiao,
    valor_formalizado,
    valor_executado,
    valor_devolvido,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__alimenta_brasil_propostas') }}
