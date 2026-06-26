# Roadmap — Open Data Agro

> Phases align with `.local/IMPLEMENTATION-PLAN.md` (detailed tasks in `.local/phases/`).  
> Pending work tracker: `.local/PENDING-TASKS.md` (gitignored).

---

## Objective

**Public agro data ingestor** — CONAB portal first, extensible to IBGE, INMET, BCB, CEPEA, ANP.

**Initial source:** [CONAB — Downloads de Arquivos](https://portaldeinformacoes.conab.gov.br/download-arquivos.html)

**Current focus:** collection sprint complete — offline pipelines + CI gates; next: doc sync and optional analytics phase (see `.local/DATA-CROSSING-VISION.md`).

---

## Stack (local-first)

| Layer | Technology | Local default |
|-------|------------|---------------|
| Ingestão / orquestração | **Go** (`cmd/ingestor`, `cmd/processor`) | `make build` |
| Data Lake | **Parquet** bronze | `./lake/bronze/` |
| Lakehouse | **Parquet** silver/gold via dbt | `./lake/silver/`, `./lake/gold/` |
| Object storage (prod) | MinIO / Cloudflare R2 | `STORAGE_MODE=local` (default) |
| Processamento | **DuckDB** | CLI + `DUCKDB_PATH` |
| Transformação | **dbt** | DuckDB profile under `dbt/` |
| Qualidade | **Great Expectations** | `expectations/` on bronze |
| DB operacional | **PostgreSQL 18** | Docker Compose |
| DB analítico | **DuckDB** | `duckdb/open_data_agro.duckdb` |

---

## Phase 0 — Platform scaffold

- [x] Agent harness + governance docs
- [x] `.local` implementation plan + phase folders
- [x] Go workspace (`go.work`), Docker Postgres + Redis
- [x] CI pipeline — [`.github/workflows/ci.yml`](../.github/workflows/ci.yml)

**Local CI mirror:** `make ci-go`, `make ci-dbt`

---

## Phase 1–8 — Platform

| Phase | Task | Status | Verify |
|-------|------|--------|--------|
| 1 | Ingestor CLI (`cmd/ingestor`) | **Done** | `go run ./cmd/ingestor run conab.estimativa-graos` |
| 2 | Parquet bronze layout | **Done** | `./lake/bronze/` |
| 3 | Silver/gold layout | **Done** | dbt → `./lake/gold/` |
| 4 | DuckDB processing + promote | **Done** | `make build-processor` |
| 5 | dbt transforms | **Done** | `make dbt-build` |
| 6 | Great Expectations gates | **Done** | `GE_INTEGRATION=1 go test ./internal/processor -run Quality` |
| 7 | PostgreSQL ops DB | **Done** | `make migrate-up` |
| 8 | DuckDB analytics views | **Done** | `make analytics-init` |

---

## Phase 10–14 — CONAB (+ ANP)

| Phase | Section | Status | MVP target |
|-------|---------|--------|------------|
| 10 | Produção Agrícola | **Done** | `make conab-mvp` — grãos, cana, café, custo produção |
| 11 | Mercado | **Done** | `make conab-mercado-full-mvp` — 8 datasets incl. prohort |
| 12 | Abastecimento + ANP | **Done** | `make conab-abastecimento-mvp`, `make anp-mvp` |
| 13 | Armazenamento e Logística | **Done** | `make conab-armazenamento-logistica-mvp` — frete + capacidade |
| 14 | Agricultura Familiar | **Done** | `make conab-agricultura-familiar-mvp` |

---

## Phase 15–19 — Multi-source collection

| Phase | Agency | Status | MVP target |
|-------|--------|--------|------------|
| 15 | IBGE Localidades | **Done** | `make ibge-localidades-mvp`, `make ibge-localidades-live-smoke` |
| 16 | IBGE PAM | **Done** | `make ibge-pam-mvp` |
| 17 | INMET Clima | **Done** | `make inmet-clima-mvp` |
| 18 | BCB SGS macro | **Done** | `make bcb-sgs-mvp` |
| 19 | CEPEA preços agro | **Done** | `make cepea-indicadores-mvp` |

**Collection sprint exit (local):** `make collection-full-mvp`  
**Collection sprint exit (CI):** `make ci-collection-full-mvp`

**Cross-cutting validation:** `make validate-codigo-ibge`, `make ci-validate-codigo-ibge`

**Refresh schedule:** [REFRESH-POLICY.md](REFRESH-POLICY.md)  
**Historical backfill:** [DATA-HISTORY-RANGES.md](DATA-HISTORY-RANGES.md)

---

## Phase 21–23 — Brazil additional official sources

| Phase | Agency | Status | MVP target |
|-------|--------|--------|------------|
| 21 | MDIC Comex Stat | **Done** | `make mdic-comex-mvp`, `make ci-mdic-comex-mvp` |
| 22 | ANTT pedágios (logística) | **Done** | `make dnit-antt-logistica-mvp`, `make ci-dnit-antt-logistica-mvp` |
| 23 | MAPA ZARC tábua de risco | **Done** | `make mapa-dados-mvp`, `make ci-mapa-dados-mvp` |
| 24 | B3 futuros agro (SOY/CCM/BGI) | **Done** | `make b3-futuros-mvp`, `make ci-b3-futuros-mvp` |
| 25 | USDA FAS PSD (soja/milho/trigo) | **Done** | `make usda-psd-mvp`, `make ci-usda-psd-mvp` |
| 26 | FAO FAOSTAT prices agro | **Done** | `make fao-faostat-mvp`, `make ci-fao-faostat-mvp` |
| 27 | World Bank Pink Sheet monthly | **Done** | `make worldbank-commodities-mvp`, `make ci-worldbank-commodities-mvp` |
| 36 | International extended (FAO prod/trade + WB ag indices) | **Done** | `make international-extended-mvp`, `make ci-international-extended-mvp` |
| 28 | NOAA climate global indices (ONI + global temp) | **Done** | `make noaa-climate-mvp`, `make ci-noaa-climate-mvp` |
| 29 | Unified PostgreSQL (gold → analytics schema) | **Done** | `make unified-db-sync`, `make ci-unified-db-sync` |

---

## Phase 33 — Collection hardening

| Task | Status | Verify |
|------|--------|--------|
| Historical ranges doc | **Done** | [DATA-HISTORY-RANGES.md](DATA-HISTORY-RANGES.md) |
| Refresh policy | **Done** | [REFRESH-POLICY.md](REFRESH-POLICY.md) |
| CONAB portal audit | **Done** | `.local/phases/33-collection-hardening/CONAB-PORTAL-AUDIT.md` |
| CI gate | **Done** | `make ci-collection-hardening-mvp` |
| Local gate | **Done** | `make collection-hardening-mvp` |

---

## CI jobs (GitHub Actions)

| Job | Purpose | Local mirror |
|-----|---------|--------------|
| `security` | `govulncheck` | — |
| `go` | unit tests, DuckDB + GE integration, lint | `make ci-go` |
| `dbt` | dbt build, analytics smoke, cod_ibge validation, full collection | `make ci-dbt` (full mirror) |
| `quality` | bronze GE checkpoint | `scripts/quality/run_checkpoint.py` |

---

## MVP path (achieved)

```text
00-platform → 01-ingestor → 02-lake → 06-GE → 07-postgres → 10-conab-mvp
→ 11-mercado → 12-abastecimento → 13-logística → 14-PAA
→ 15-localidades → 16-PAM → 17-INMET → 18-BCB → 19-CEPEA
→ collection-full-mvp (sprint exit)
```

---

## Next (post-collection)

**Ingestor expansion (`.local/`):**

- Phase 32 — source discovery ✅ — [SOURCE-DISCOVERY-CATALOG.md](../.local/SOURCE-DISCOVERY-CATALOG.md) (local, gitignored)
- Phase 33 — collection hardening ✅ — [DATA-HISTORY-RANGES.md](DATA-HISTORY-RANGES.md) · `make ci-collection-hardening-mvp`
- Phase 34 — BR logistics extended ✅ — ANTT volume + receita · `make ci-br-logistica-extended-mvp`
- Phase 35 — BR comex extended ✅ — MDIC import + export UF + diesel · `make ci-mdic-comex-extended-mvp`
- Phase 36 — international extended ✅ — FAO production/trade + World Bank ag indices · `make ci-international-extended-mvp`
- Phases 37–38 — gap closure (IBGE LSPA, EIA)
- Phase 29 — unified PostgreSQL ✅ — `make unified-db-sync`

**Phase 20 scaffold:** [.local/phases/20-analytics-crossing/README.md](../.local/phases/20-analytics-crossing/README.md) — analytics crossing (feature joins); implementation **not started** (IA deferred).

Remaining tracks (see `.local/PENDING-TASKS.md`):

- Optional: GE bronze referential `cod_ibge` (deferred — post-dbt script is gate)

**Phase 2–3 storage:** `make ci-minio` · `make ci-validate-r2-env` · `make ci-delta-versioning` (native Delta silver versioning)

---

## Related

- [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md) — per-dataset catalog
- [NEW-PROJECT-CHECKLIST.md](NEW-PROJECT-CHECKLIST.md)
- `.local/IMPLEMENTATION-PLAN.md` (gitignored)
- `.local/PENDING-TASKS.md` (gitignored)
