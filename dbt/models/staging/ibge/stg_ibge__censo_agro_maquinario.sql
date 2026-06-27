with silver as (
    select * from {{ read_silver_parquet('censo_agro_maquinario', agency='ibge') }}
),

renamed as (
    select
        {{ ibge_censo_agro_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        'https://censoagro2017.ibge.gov.br/' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by codigo_uf, ano, variavel_codigo, condicao_produtor_codigo, tipologia_codigo, atividade_codigo
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
