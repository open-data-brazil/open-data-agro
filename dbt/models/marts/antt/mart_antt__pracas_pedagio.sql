{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_antt__pracas_pedagio/mart.parquet',
    format='parquet'
) }}

select
    concessionaria,
    praca_de_pedagio,
    ano_do_pnv_snv,
    rodovia,
    uf,
    km_m,
    municipio,
    tipo_de_pista,
    sentido,
    situacao,
    data_da_inativacao,
    latitude,
    longitude,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_antt__pracas_pedagio') }}
