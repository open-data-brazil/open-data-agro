with silver as (
    select * from {{ read_silver_parquet('global_temp_anomaly', agency='noaa') }}
),

renamed as (
    select
        refmonth,
        anomaly,
        unit,
        base_period,
        index_name,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ noaa_fonte_oficial("noaa.global-temp-anomaly") }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by refmonth, index_name
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
