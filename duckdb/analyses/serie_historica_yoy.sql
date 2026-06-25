-- Year-over-year production change from historical series (local silver view).
-- Filter produto/uf to prune scan volume on laptop hardware.

SELECT
    produto,
    uf,
    ano,
    try_cast(producao_mil_t AS DOUBLE) AS producao_mil_t,
    try_cast(producao_mil_t AS DOUBLE)
        - lag(try_cast(producao_mil_t AS DOUBLE)) OVER (
            PARTITION BY produto, uf
            ORDER BY try_cast(ano AS INTEGER)
        ) AS yoy_change_mil_t
FROM analytics.conab_serie_historica_graos
WHERE produto = 'Soja'
  AND uf IN ('PR', 'RS', 'MT')
ORDER BY uf, try_cast(ano AS INTEGER);
