with silver as (
    select * from {{ read_silver_parquet('pluviometria_redes', agency='ana') }}
),

renamed as (
    select
        station_code,
        consistency_level,
        data_type,
        observed_at,
        daily_mean,
        acquisition_method,
        max_value,
        min_value,
        mean_value,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ ana_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by station_code, observed_at, data_type, consistency_level
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
