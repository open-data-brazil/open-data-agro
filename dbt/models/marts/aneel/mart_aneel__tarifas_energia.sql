{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_aneel__tarifas_energia/mart.parquet',
    format='parquet'
) }}

select
    DatGeracaoConjuntoDados,
    DatCompetencia,
    NomBandeiraAcionada,
    VlrAdicionalBandeira,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_aneel__tarifas_energia') }}
