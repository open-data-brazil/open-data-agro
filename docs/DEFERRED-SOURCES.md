# Deferred sources — unreachable without proxy, broken URLs, or fixture-only

These datasets were **removed from the active catalog** (`configs/catalog/`) on 2026-06-27 after manual validation (curl + ingestor `DownloadSource`). They are kept in `configs/catalog/_deferred/unreachable_sources.yaml` for reference only — the loader **skips** `_deferred/` paths, so neither the ingestor nor the source-health bot probes them.

## Summary

| Dataset ID | Why deferred | Live access path (future) |
|------------|--------------|---------------------------|
| `antaq.movimentacao-carga-portuaria` | Bulk URL returns **HTTP 404** (deprecated Painel Estatístico Aquaviário export) | Find new ANTAQ open-data bulk/API on [gov.br/antaq dados abertos](https://www.gov.br/antaq/pt-br/acesso-a-informacao/dados-abertos); update `source_url`; move entry back to `configs/catalog/antaq/` |
| `fao.comercio-agro` | FAOSTAT **Trade** bulk ZIP returns **HTTP 403** (Production/Prices ZIPs still work) | Monitor [FAOSTAT TCL bulk](https://www.fao.org/faostat/en/#data/TCL); when `Trade_Crops_Livestock_E_All_Data_(Normalized).zip` is publicly reachable again, restore catalog entry |
| `mexico.siap-produccion-agricola` | No stable machine-readable API; ingest used **embedded fixture only** | Obtain SIAP CKAN/API bulk from [gob.mx/siap](https://www.gob.mx/siap) or set `SIAP_BULK_PATH` to a verified JSON export; implement real HTTP fetch in `internal/mexico/` |
| `noaa.gpcc-precipitation` | Catalog URLs **404**; ingest used **embedded fixture** | Use stable GPCC monthly grid from [NOAA NCEI GPCC](https://www.ncei.noaa.gov/products/land-based-station/global-precipitation-climatology-centre) or DWD opendata; update `source_url` and remove fixture path in `internal/noaa/gpcc.go` |
| `usda.gats-trade` | `apps.fas.usda.gov` **connection refused** from BR/CI networks + requires `USDA_FAS_API_KEY` | Run ingest from US egress (or when USDA FAS is reachable); register at [USDA FAS Open Data](https://apps.fas.usda.gov/opendata/); set `USDA_FAS_API_KEY` |
| `usda.psd-soja` | Same host **connection refused**; SOAP unreachable without US network | Run from US egress; verify SOAP `getDatabyCommodityPerYear` on `PSDExternalAPIService/svcPSD_AMIS.asmx` |
| `usda.psd-milho` | Same as PSD soja | Same as above |
| `usda.psd-trigo` | Same as PSD soja | Same as above |
| `wto.its-trade-statistics` | API returns **401** without key; ingest used **fixture** without `WTO_API_KEY` | Subscribe at [WTO API portal](https://api.wto.org/); set repo secret `WTO_API_KEY`; remove fixture fallback or gate on key presence |

## Re-enable checklist

1. Confirm live download with curl or `go test -tags=integration` against the agency client.
2. Move the YAML block from `configs/catalog/_deferred/unreachable_sources.yaml` into a new file under `configs/catalog/<agency>/`.
3. Restore dbt staging/mart, DuckDB view, GE suite, and Makefile MVP target (see git history before 2026-06-27).
4. Run the relevant `make *-mvp` gate and source-health bot locally.

## Related — still active but need official API keys

These remain in the catalog; warnings are expected until secrets are configured:

| Dataset ID | Secret | Registration |
|------------|--------|--------------|
| `eia.petroleum-prices` | `EIA_API_KEY` | [EIA Open Data](https://www.eia.gov/opendata/register.php) |

## Catalog loader behavior

`LoadRegistryDir` skips any path containing a `_deferred` directory segment. Deferred entries are **not** counted in `docs/SOURCE-HEALTH.md` daily probes.
