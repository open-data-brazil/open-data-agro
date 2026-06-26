{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_mapa__zarc_tabua_risco/mart.parquet',
    format='parquet'
) }}

select
    nome_cultura,
    safra_ini,
    safra_fim,
    cod_cultura,
    cod_ciclo,
    cod_solo,
    geocodigo,
    uf,
    municipio,
    cod_clima,
    nome_clima,
    cod_outros_manejos,
    nome_outros_manejos,
    produtividade,
    cod_nm,
    cod_munic,
    cod_meso,
    cod_micro,
    portaria,
    dec1, dec2, dec3, dec4, dec5, dec6, dec7, dec8, dec9, dec10,
    dec11, dec12, dec13, dec14, dec15, dec16, dec17, dec18, dec19, dec20,
    dec21, dec22, dec23, dec24, dec25, dec26, dec27, dec28, dec29, dec30,
    dec31, dec32, dec33, dec34, dec35, dec36,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_mapa__zarc_tabua_risco') }}
