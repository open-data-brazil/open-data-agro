{% macro wto_its_trade_statistics_fonte_oficial() %}
  https://stats.wto.org/
{% endmacro %}

{% macro wto_its_trade_statistics_columns() %}
    reporter_code, reporter_name, partner_code, partner_name, indicator_code, period, value_usd, flow_code
{% endmacro %}

{% macro wto_its_trade_statistics_dedupe_keys() %}
    reporter_code, indicator_code, period, flow_code
{% endmacro %}
