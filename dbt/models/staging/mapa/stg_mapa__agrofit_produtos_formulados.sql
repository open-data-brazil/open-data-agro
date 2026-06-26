with silver as (
    select * from {{ read_silver_parquet('agrofit_produtos_formulados', agency='mapa') }}
),

renamed as (
    select
        NR_REGISTRO as nr_registro,
        MARCA_COMERCIAL as marca_comercial,
        FORMULACAO as formulacao,
        INGREDIENTE_ATIVO as ingrediente_ativo,
        TITULAR_DE_REGISTRO as titular_de_registro,
        CLASSE as classe,
        MODO_DE_ACAO as modo_de_acao,
        CULTURA as cultura,
        PRAGA_NOME_CIENTIFICO as praga_nome_cientifico,
        PRAGA_NOME_COMUM as praga_nome_comum,
        EMPRESA_PAIS_TIPO as empresa_pais_tipo,
        CLASSE_TOXICOLOGICA as classe_toxicologica,
        CLASSE_AMBIENTAL as classe_ambiental,
        ORGANICOS as organicos,
        SITUACAO as situacao,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ mapa_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by nr_registro, marca_comercial, cultura, praga_nome_cientifico
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
