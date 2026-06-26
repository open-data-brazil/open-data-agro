# Refresh policy — Open Data Agro ingestor

> Public schedule for refreshing official datasets.  
> **Scope:** ingestor maintenance only — not IA retraining.

**Local-first:** all targets runnable via `make` without cloud. Production cron uses same commands with `STORAGE_MODE=r2` when configured.

---

## Schedule summary

| Agency | Datasets | Frequency | Local target | Notes |
|--------|----------|-----------|--------------|-------|
| **CONAB** | Produção, Mercado, Logística, Abastecimento, PAA | Weekly (Mon) | `make conab-mercado-full-mvp` (subset) | Portal updates vary by section |
| **ANP** | LPC combustíveis | Weekly | `make anp-mvp` | ANP publishes weekly |
| **BCB** | SGS IPCA, PTAX, IGP-M | Daily (business days) | `make bcb-sgs-mvp` | API live fetch |
| **CEPEA** | Indicadores agro | Daily | `make cepea-indicadores-mvp` | Portal or mirror fallback |
| **IBGE** | Localidades | Quarterly | `make ibge-localidades-mvp` | Rarely changes |
| **IBGE** | PAM | Annual (after release) | `make ibge-pam-mvp` | Post-PAM publication |
| **INMET** | BDMEP annual ZIP | Annual (Jan) | `make inmet-clima-mvp` | New year ZIP each January |
| **MDIC** | Comex export | Monthly (after Secex) | `make mdic-comex-mvp` | ~30 days after month close |
| **ANTT** | Pedágios | Monthly | `make dnit-antt-logistica-mvp` | CKAN CSV update |
| **MAPA** | ZARC | Annual (safra) | `make mapa-dados-mvp` | New safra CSV |
| **B3** | Futuros agro | Daily (business days) | `make b3-futuros-mvp` | SPRD zip per pregão |
| **USDA** | PSD | Monthly | `make usda-psd-mvp` | AMIS web service |
| **FAO** | FAOSTAT prices | Annual / on bulk release | `make fao-faostat-mvp` | Bulk zip refresh |
| **World Bank** | Pink Sheet | Monthly | `make worldbank-commodities-mvp` | CMO workbook |
| **NOAA** | Climate indices | Monthly | `make noaa-climate-mvp` | ONI + global temp CSV |

---

## Full collection (sprint / CI)

| Target | When | Purpose |
|--------|------|---------|
| `make collection-full-mvp` | Manual sprint exit | All BR + intl core E2E local |
| `make ci-collection-full-mvp` | CI on every PR | Regression gate |
| `make unified-db-sync` | After gold marts update | PostgreSQL `analytics.*` refresh |

---

## Historical backfill

First-time or recovery backfill uses the same ingestor with explicit `--from`:

```bash
go run ./cmd/ingestor run bcb.sgs-ptax-usd-venda --from 1984-01-01
go run ./cmd/ingestor run mdic.comex-exportacao-ncm-mes --from 1997-01-01
```

Document per-dataset earliest date in [DATA-HISTORY-RANGES.md](DATA-HISTORY-RANGES.md) and phase `OFFICIAL-REFERENCE.md` files (Phase 33).

---

## Production (R2)

When `STORAGE_MODE=r2`:

1. Same `cmd/ingestor run` commands — object keys unchanged
2. Schedule via cron / GitHub Actions / external orchestrator
3. Validate env: `make validate-r2-env`

---

## Related

- [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md)
- [DATA-HISTORY-RANGES.md](DATA-HISTORY-RANGES.md) — per-dataset backfill `--from` and source limits (Phase 33)
- [POSTGRES-UNIFIED-SYNC.md](POSTGRES-UNIFIED-SYNC.md)
- `.local/phases/33-collection-hardening/`
