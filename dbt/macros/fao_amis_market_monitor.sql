{% macro fao_amis_market_monitor_fonte_oficial() %}
  https://www.amis-outlook.org/
{% endmacro %}

{% macro fao_amis_market_monitor_columns() %}
    commodity_slug, refmonth, indicator_slug, value, unit
{% endmacro %}

{% macro fao_amis_market_monitor_dedupe_keys() %}
    commodity_slug, refmonth, indicator_slug
{% endmacro %}
