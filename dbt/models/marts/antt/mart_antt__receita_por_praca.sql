{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_antt__receita_por_praca/mart.parquet',
    format='parquet'
) }}

select
    concessionaria,
    praca_de_pedagio,
    ano_pnv_snv,
    uf,
    rodovia,
    km_m,
    tipo_de_pista,
    sentido,
    municipio,
    direcao,
    latitude,
    longitude,
    data_da_ativacao,
    mes_ano,
    receita_praca_de_pedagio,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_antt__receita_por_praca') }}
