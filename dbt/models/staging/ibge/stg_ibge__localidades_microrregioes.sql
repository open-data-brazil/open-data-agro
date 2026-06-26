with silver as (
    select * from {{ read_silver_parquet('localidades_microrregioes', agency='ibge') }}
),

renamed as (
    select
        {{ ibge_localidades_microrregioes_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ ibge_localidades_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by codigo_microrregiao
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
