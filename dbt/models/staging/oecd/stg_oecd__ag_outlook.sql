with silver as (
    select * from {{ read_silver_parquet('ag_outlook', agency='oecd') }}
),

renamed as (
    select
        {{ oecd_ag_outlook_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ oecd_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by commodity_code, measure_code, year
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
