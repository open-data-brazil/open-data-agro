with silver as (
    select * from {{ read_silver_parquet('deter_alertas_desmatamento', agency='inpe') }}
),

renamed as (
    select
        *,
        cast(_ingested_at as varchar) as capturado_em,
        'https://terrabrasilis.dpi.inpe.br/downloads/' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by view_date
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
