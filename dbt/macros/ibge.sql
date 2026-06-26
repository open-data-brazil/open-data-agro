{% macro ibge_fonte_oficial() %}
  https://sidra.ibge.gov.br/pesquisa/pam
{% endmacro %}

{% macro ibge_localidades_fonte_oficial() %}
  https://servicodados.ibge.gov.br/api/docs/localidades
{% endmacro %}

{% macro ibge_localidades_municipios_columns() %}
    codigo_ibge,
    nome,
    sigla_uf,
    codigo_uf,
    codigo_regiao,
    nome_regiao
{% endmacro %}

{% macro ibge_localidades_ufs_columns() %}
    codigo_uf,
    sigla_uf,
    nome,
    codigo_regiao,
    sigla_regiao,
    nome_regiao
{% endmacro %}

{% macro ibge_pam_area_quantidade_columns() %}
    sidra_tabela,
    codigo_ibge,
    municipio,
    ano,
    variavel_codigo,
    variavel,
    produto_codigo,
    produto,
    valor,
    unidade_codigo,
    unidade
{% endmacro %}
