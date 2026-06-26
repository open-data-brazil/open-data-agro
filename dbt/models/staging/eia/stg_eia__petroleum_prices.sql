with silver as (
    select * from {{ read_silver_parquet('petroleum_prices', agency='eia') }}
),

renamed as (
    select
        {{ eia_petroleum_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ eia_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by series_id, refdate
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
