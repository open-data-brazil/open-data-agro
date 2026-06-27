with silver as (
    select * from {{ read_silver_parquet('carga_energetica', agency='ons') }}
),

renamed as (
    select
        *,
        cast(_ingested_at as varchar) as capturado_em,
        'https://dados.ons.org.br/dataset/carga-energia' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by id_subsistema
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
