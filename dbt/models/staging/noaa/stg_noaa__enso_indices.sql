with silver as (
    select * from {{ read_silver_parquet('enso_indices', agency='noaa') }}
),

renamed as (
    select
        refyear,
        season_code,
        sst_total,
        anomaly,
        index_name,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ noaa_fonte_oficial("noaa.enso-indices") }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by refyear, season_code
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
