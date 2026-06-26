with silver as (
    select * from {{ read_silver_parquet('sequia_monitor', agency='inmet') }}
),

renamed as (
    select
        {{ inmet_sequia_monitor_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ inmet_sequia_monitor_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by mapa_id, categoria_seca, area_id, tipo_area
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
