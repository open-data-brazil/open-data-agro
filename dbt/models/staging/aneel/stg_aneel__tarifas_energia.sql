with silver as (
    select * from {{ read_silver_parquet('tarifas_energia', agency='aneel') }}
),

renamed as (
    select
        {{ aneel_tarifas_energia_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ aneel_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by DatCompetencia, NomBandeiraAcionada
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
