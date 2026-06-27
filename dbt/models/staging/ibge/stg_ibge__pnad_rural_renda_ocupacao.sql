with silver as (
    select * from {{ read_silver_parquet('pnad_rural_renda_ocupacao', agency='ibge') }}
),

renamed as (
    select
        {{ ibge_pnad_rural_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        'https://sidra.ibge.gov.br/pesquisa/pnad' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by codigo_uf, trimestre, variavel_codigo
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
