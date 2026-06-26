{% macro fao_ffpi_fonte_oficial() %}
  https://www.fao.org/worldfoodsituation/foodpricesindex/en/
{% endmacro %}

{% macro fao_ffpi_columns() %}
    refmonth,
    index_slug,
    index_name,
    value,
    base_period
{% endmacro %}
