{% macro anp_silver_path(table_name) %}
  {{ var('lake_root') }}/silver/anp/{{ table_name }}
{% endmacro %}

{% macro anp_fonte_oficial() %}
  {{ return(var('anp_fonte_oficial')) }}
{% endmacro %}

{% macro anp_combustiveis_precos_medios_columns() %}
  trim("DATA INICIAL") as data_inicial,
  trim("DATA FINAL") as data_final,
  trim("ESTADO") as estado,
  trim("MUNICÍPIO") as municipio,
  trim("PRODUTO") as produto,
  cast("NÚMERO DE POSTOS PESQUISADOS" as varchar) as numero_postos_pesquisados,
  trim("UNIDADE DE MEDIDA") as unidade_medida,
  cast("PREÇO MÉDIO REVENDA" as varchar) as preco_medio_revenda,
  cast("DESVIO PADRÃO REVENDA" as varchar) as desvio_padrao_revenda,
  cast("PREÇO MÍNIMO REVENDA" as varchar) as preco_minimo_revenda,
  cast("PREÇO MÁXIMO REVENDA" as varchar) as preco_maximo_revenda,
  cast("COEF DE VARIAÇÃO REVENDA" as varchar) as coef_variacao_revenda
{% endmacro %}

{% macro anp_combustiveis_precos_postos_columns() %}
  trim("CNPJ") as cnpj,
  trim("RAZÃO") as razao,
  trim("FANTASIA") as fantasia,
  trim("ENDEREÇO") as endereco,
  trim("NÚMERO") as numero,
  trim("COMPLEMENTO") as complemento,
  trim("BAIRRO") as bairro,
  trim("CEP") as cep,
  trim("MUNICÍPIO") as municipio,
  trim("ESTADO") as estado,
  trim("BANDEIRA") as bandeira,
  trim("PRODUTO") as produto,
  trim("UNIDADE DE MEDIDA") as unidade_medida,
  cast("PREÇO DE REVENDA" as varchar) as preco_revenda,
  trim("DATA DA COLETA") as data_coleta
{% endmacro %}
