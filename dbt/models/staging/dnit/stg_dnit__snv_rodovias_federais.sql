with silver as (
    select * from {{ read_silver_parquet('snv_rodovias_federais', agency='dnit') }}
),

renamed as (
    select
        {{ dnit_snv_rodovias_federais_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ dnit_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by br, uf, codigo, km_inicial, km_final
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
