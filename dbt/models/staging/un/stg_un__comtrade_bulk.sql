with silver as (
    select * from {{ read_silver_parquet('comtrade_bulk', agency='un') }}
),

renamed as (
    select
        {{ un_comtrade_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ un_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by reporter_code, partner_code, flow_code, period, hs_code
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
