{% macro fao_fonte_oficial() %}
  https://www.fao.org/faostat/en/#data/PP
{% endmacro %}

{% macro fao_producao_fonte_oficial() %}
  https://www.fao.org/faostat/en/#data/QCL
{% endmacro %}

{% macro fao_comercio_fonte_oficial() %}
  https://www.fao.org/faostat/en/#data/TCL
{% endmacro %}

{% macro fao_annual_bulk_columns() %}
    area_code,
    area_name,
    item_code,
    item_name,
    commodity_slug,
    element_code,
    element_name,
    year,
    unit,
    value,
    flag
{% endmacro %}
