with silver as (
    select * from {{ read_silver_parquet('receita_por_praca', agency='antt') }}
),

renamed as (
    select
        "Concessionaria" as concessionaria,
        "Praca_de_pedagio" as praca_de_pedagio,
        "Ano_PNV_SNV" as ano_pnv_snv,
        "UF" as uf,
        "Rodovia" as rodovia,
        "Km_m" as km_m,
        "Tipo_de_Pista" as tipo_de_pista,
        "Sentido" as sentido,
        "Municipio" as municipio,
        "Direcao" as direcao,
        "Latitude" as latitude,
        "Longitude" as longitude,
        "Data_da_Ativacao" as data_da_ativacao,
        "Mes_ano" as mes_ano,
        "Receita_Praca_de_Pedagio" as receita_praca_de_pedagio,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ antt_receita_por_praca_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * from renamed
