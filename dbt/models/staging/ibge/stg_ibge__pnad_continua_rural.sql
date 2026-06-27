with silver as (
    select * from {{ read_silver_parquet('pnad_continua_rural', agency='ibge') }}
),

renamed as (
    select
        *,
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
            partition by codigo_uf
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
