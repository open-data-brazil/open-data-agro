with silver as (
    select * from {{ read_silver_parquet('commodity_indexes', agency='fred') }}
),

renamed as (
    select
        {{ fred_commodity_indexes_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ fred_commodity_indexes_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by series_id, refmonth
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
