with silver as (
    select * from {{ read_silver_parquet('goi_index', agency='igc') }}
),

renamed as (
    select
        {{ igc_goi_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ igc_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by index_slug, refdate
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
