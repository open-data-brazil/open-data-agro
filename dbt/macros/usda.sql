{% macro usda_fonte_oficial() %}
  https://www.usda.gov/oce/commodity-markets/wasde
{% endmacro %}

{% macro usda_wasde_columns() %}
    report_month,
    commodity,
    market_year,
    attribute,
    value,
    unit
{% endmacro %}
