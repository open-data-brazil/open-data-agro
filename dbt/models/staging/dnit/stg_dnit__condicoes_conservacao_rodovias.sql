with silver as (
    select * from {{ read_silver_parquet('condicoes_conservacao_rodovias', agency='dnit') }}
),

renamed as (
    select
        *,
        cast(_ingested_at as varchar) as capturado_em,
        'https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by id_malha
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
