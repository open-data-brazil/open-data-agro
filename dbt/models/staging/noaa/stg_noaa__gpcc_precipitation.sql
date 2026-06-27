with silver as (
    select * from {{ read_silver_parquet('gpcc_precipitation', agency='noaa') }}
),

renamed as (
    select
        {{ noaa_gpcc_precipitation_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ noaa_gpcc_precipitation_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by refmonth, latitude, longitude
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
