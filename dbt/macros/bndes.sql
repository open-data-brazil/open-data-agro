{% macro bndes_fonte_oficial() %}
  https://dadosabertos.bndes.gov.br/dataset/desembolsos
{% endmacro %}

{% macro bndes_financiamento_agro_columns() %}
    ano,
    mes,
    agropecuaria
{% endmacro %}

{% macro bndes_desembolsos_linhas_agro_columns() %}
    ano,
    mes,
    bndes_finem,
    bndes_exim,
    bndes_mercado_de_capitais,
    bndes_nao_reembolsavel,
    bndes_microcredito,
    bndes_prestacao_de_garantia,
    bndes_finame
{% endmacro %}
