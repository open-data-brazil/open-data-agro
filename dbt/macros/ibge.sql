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

{% macro ibge_localidades_regioes_columns() %}
    codigo_regiao,
    sigla_regiao,
    nome
{% endmacro %}

{% macro ibge_localidades_mesorregioes_columns() %}
    codigo_mesorregiao,
    nome,
    codigo_uf,
    sigla_uf,
    nome_uf,
    codigo_regiao,
    sigla_regiao,
    nome_regiao
{% endmacro %}

{% macro ibge_localidades_microrregioes_columns() %}
    codigo_microrregiao,
    nome,
    codigo_mesorregiao,
    nome_mesorregiao,
    codigo_uf,
    sigla_uf,
    nome_uf
{% endmacro %}

{% macro ibge_lspa_fonte_oficial() %}
  https://sidra.ibge.gov.br/pesquisa/lspa
{% endmacro %}

{% macro ibge_lspa_area_producao_columns() %}
    sidra_tabela,
    codigo_uf,
    uf,
    mes,
    variavel_codigo,
    variavel,
    produto_codigo,
    produto,
    produto_slug,
    valor,
    unidade_codigo,
    unidade
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

{% macro ibge_pam_core_columns() %}
    {{ ibge_pam_area_quantidade_columns() }}
{% endmacro %}

{% macro ibge_pevs_fonte_oficial() %}
  https://sidra.ibge.gov.br/pesquisa/pevs
{% endmacro %}

{% macro ibge_pevs_producao_vegetal_columns() %}
    sidra_tabela,
    codigo_uf,
    uf,
    ano,
    variavel_codigo,
    variavel,
    produto_codigo,
    produto,
    valor,
    unidade_codigo,
    unidade
{% endmacro %}

{% macro ibge_ppm_fonte_oficial() %}
  https://sidra.ibge.gov.br/pesquisa/pam
{% endmacro %}

{% macro ibge_ppm_producao_municipal_columns() %}
    {{ ibge_pam_core_columns() }}
{% endmacro %}

{% macro ibge_ppm_uf_columns() %}
    sidra_tabela,
    codigo_uf,
    uf,
    ano,
    variavel_codigo,
    variavel,
    categoria_codigo,
    categoria,
    valor,
    unidade_codigo,
    unidade
{% endmacro %}

{% macro ibge_censo_agro_columns() %}
    sidra_tabela,
    codigo_uf,
    uf,
    ano,
    variavel_codigo,
    variavel,
    condicao_produtor_codigo,
    condicao_produtor,
    tipologia_codigo,
    tipologia,
    atividade_codigo,
    atividade,
    sexo_produtor_codigo,
    sexo_produtor,
    idade_produtor_codigo,
    idade_produtor,
    valor,
    unidade_codigo,
    unidade
{% endmacro %}

{% macro ibge_pnad_rural_columns() %}
    sidra_tabela,
    codigo_uf,
    uf,
    trimestre,
    variavel_codigo,
    variavel,
    valor,
    unidade_codigo,
    unidade
{% endmacro %}
