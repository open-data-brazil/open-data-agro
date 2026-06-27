{% macro cftc_cot_agricultural_futures_fonte_oficial() %}
  https://www.cftc.gov/MarketReports/CommitmentsofTraders/index.htm
{% endmacro %}

{% macro cftc_cot_agricultural_futures_columns() %}
    report_date, commodity_name, commodity_slug, market_name, open_interest_all, m_money_long, m_money_short, prod_merc_long, prod_merc_short, commodity_group, futonly_or_combined
{% endmacro %}

{% macro cftc_cot_agricultural_futures_dedupe_keys() %}
    commodity_slug, report_date
{% endmacro %}
