with silver as (
    select * from {{ read_silver_parquet('agroapi_agrofit', agency='embrapa') }}
),

renamed as (
    select
        *,
        cast(_ingested_at as varchar) as capturado_em,
        'https://www.agroapi.cnptia.embrapa.br/store/' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by _dataset_id, numero_registro, marca_comercial
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
