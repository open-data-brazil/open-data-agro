-- Total estimated production by crop and UF (latest season in mart).
-- Local-only: requires make analytics-init after dbt gold build.

SELECT
    produto,
    uf,
    safra,
    sum(try_cast(producao_mil_t AS DOUBLE)) AS producao_mil_t_total
FROM analytics.conab_estimativa_graos
WHERE safra = (
    SELECT max(safra)
    FROM analytics.conab_estimativa_graos
)
GROUP BY produto, uf, safra
ORDER BY produto, uf;
