with silver as (
    select * from {{ read_silver_parquet('comercio_mercadorias_zfm', agency='suframa') }}
),

renamed as (
    select
        *,
        cast(_ingested_at as varchar) as capturado_em,
        'https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by column_1
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
