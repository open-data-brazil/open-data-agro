{% macro ibge_fonte_oficial() %}
  https://sidra.ibge.gov.br/pesquisa/pam
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
