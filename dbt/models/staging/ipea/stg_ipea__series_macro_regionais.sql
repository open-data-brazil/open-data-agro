with silver as (
    select * from {{ read_silver_parquet('series_macro_regionais', agency='ipea') }}
),

renamed as (
    select
        {{ ipea_series_macro_regionais_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ ipea_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by series_code, refdate, territory_code
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
