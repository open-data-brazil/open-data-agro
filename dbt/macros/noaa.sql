{% macro noaa_fonte_oficial(dataset_id) %}
  {% if dataset_id == 'noaa.enso-indices' %}
    https://www.cpc.ncep.noaa.gov/products/analysis_monitoring/ensostuff/ONI_v5.php
  {% else %}
    https://www.ncei.noaa.gov/access/monitoring/climate-at-a-glance/global/time-series
  {% endif %}
{% endmacro %}
