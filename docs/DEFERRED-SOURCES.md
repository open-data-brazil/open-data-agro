# Deferred sources — unreachable without proxy, broken URLs, or fixture-only

These datasets were **removed from the active catalog** (`configs/catalog/`) on 2026-06-27 after manual validation (curl + ingestor `DownloadSource`). They are kept in `configs/catalog/_deferred/unreachable_sources.yaml` for reference only — the loader **skips** `_deferred/` paths, so neither the ingestor nor the source-health bot probes them.

**Phase 51 re-audit (2026-06-27):** all nine entries below were probed again from CI network. **None** are ready for re-enable; see [`.local/phases/51-source-discovery-wave-5/DISCOVERY-REPORT-WAVE5.md`](../.local/phases/51-source-discovery-wave-5/DISCOVERY-REPORT-WAVE5.md) (local, gitignored).

## Summary

| Dataset ID | Why deferred | Live access path (future) | Phase 51 status |
|------------|--------------|---------------------------|-----------------|
| `antaq.movimentacao-carga-portuaria` | Bulk URL returns **HTTP 404** (deprecated Painel Estatístico Aquaviário export) | Find new ANTAQ open-data bulk/API on [gov.br/antaq dados abertos](https://www.gov.br/antaq/pt-br/acesso-a-informacao/dados-abertos); update `source_url`; move entry back to `configs/catalog/antaq/` | **Still 404** — no replacement bulk found |
| `fao.comercio-agro` | FAOSTAT **Trade** bulk ZIP returns **HTTP 403** (Production/Prices ZIPs still work) | Monitor [FAOSTAT TCL bulk](https://www.fao.org/faostat/en/#data/TCL); when `Trade_Crops_Livestock_E_All_Data_(Normalized).zip` is publicly reachable again, restore catalog entry | **Still 403** |
| `mexico.siap-produccion-agricola` | No stable machine-readable API; ingest used **embedded fixture only** | Obtain SIAP CKAN/API bulk from [gob.mx/siap](https://www.gob.mx/siap) or set `SIAP_BULK_PATH` to a verified JSON export; implement real HTTP fetch in `internal/mexico/` | Portal **200**; no bulk API |
| `noaa.gpcc-precipitation` | Catalog URLs **404**; ingest used **embedded fixture** | Use stable GPCC monthly grid from [NOAA NCEI GPCC metadata](https://www.ncei.noaa.gov/access/metadata/landing-page/bin/iso?id=gov.noaa.ncdc:C00879) or DWD opendata; update `source_url` and remove fixture path in `internal/noaa/gpcc.go` | NCEI metadata **200**; grid paths **404** |
| `usda.gats-trade` | Requires `USDA_FAS_API_KEY`; **HTTP 403** without key from BR CI | Register at [USDA FAS Open Data](https://apps.fas.usda.gov/opendata/); set `USDA_FAS_API_KEY` | Host reachable (**403** not connection refused) |
| `usda.psd-soja` | SOAP/WSDL **HTTP 404** from BR CI; US egress likely required | Run from US egress; verify SOAP `getDatabyCommodityPerYear` on `PSDExternalAPIService/svcPSD_AMIS.asmx` | **Still blocked** |
| `usda.psd-milho` | Same as PSD soja | Same as above | **Still blocked** |
| `usda.psd-trigo` | Same as PSD soja | Same as above | **Still blocked** |
| `wto.its-trade-statistics` | API returns **401** without key; ingest used **fixture** without `WTO_API_KEY` | Subscribe at [WTO API portal](https://api.wto.org/); set repo secret `WTO_API_KEY`; remove fixture fallback or gate on key presence | **Still 401** |

## Discovery deferrals (not in `_deferred/` YAML)

These remain **out of catalog** from prior waves; Phase 51 re-audit:

| Dataset ID | Blocker | Phase 51 |
|------------|---------|----------|
| `antt.tarifas-pedagio` | ANTT CKAN `praca-de-pedagio` primary resource is Power BI embed — no CSV tariff table | deferred |
| `inpe.cptec-indices-climaticos` | CPTEC FTP seasonal bulk paths HTTP 404 | deferred |
| `embrapa.solos-brasil` | Embrapa solos portal 200; stable layer ZIP URL not found | deferred |
| `mapa.sigef-areas` | **Re-approved** — SIGEF CKAN CSV verified; routes to Phase 52 (removed from deferral) | approved |

## Re-enable checklist

1. Confirm live download with curl or `python3 scripts/ci/verify_wave5_discovery_probe.py` (deferred blocker section).
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

## Validation

```bash
python3 scripts/ci/verify_wave5_discovery_probe.py
```

Expect 18/18 approved primary URLs OK and deferred blockers still non-200.
