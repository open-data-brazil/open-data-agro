{% macro dnit_fonte_oficial() %}
  https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
{% endmacro %}

{% macro dnit_snv_rodovias_federais_columns() %}
    BR as br,
    UF as uf,
    "Tipo de trecho" as tipo_trecho,
    "Código" as codigo,
    "Local de Início" as local_inicio,
    "Local de Fim" as local_fim,
    "km inicial" as km_inicial,
    "km final" as km_final,
    "Extensão" as extensao,
    "Superfície Federal" as superficie_federal,
    Obras as obras,
    "Federal Coincidente" as federal_coincidente,
    "Administração" as administracao,
    "Ato legal" as ato_legal,
    "Estadual Coincidente" as estadual_coincidente,
    "Superfície Est. Coincidente" as superficie_est_coincidente,
    "Jurisdição" as jurisdicao,
    "Superfície" as superficie,
    "Unidade Local" as unidade_local
{% endmacro %}
