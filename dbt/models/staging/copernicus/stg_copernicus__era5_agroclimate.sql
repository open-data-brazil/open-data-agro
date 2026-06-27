with silver as (
    select * from {{ read_silver_parquet('era5_agroclimate', agency='copernicus') }}
),

renamed as (
    select
        {{ copernicus_era5_agroclimate_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ copernicus_era5_agroclimate_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by latitude, longitude, refdate, variable_slug
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
