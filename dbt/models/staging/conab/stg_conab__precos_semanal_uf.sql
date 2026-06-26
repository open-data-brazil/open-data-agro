with silver as (
    select * from {{ read_silver_parquet('precos_agropecuarios_semanal_uf') }}
),

renamed as (
    select
        {{ conab_precos_semanal_uf_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ conab_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by produto, id_produto, uf, ano, mes, semana, nivel_comercializacao
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
