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

---

## Phase 21–23 — Brazil additional official sources

| Phase | Agency | Status | MVP target |
|-------|--------|--------|------------|
| 21 | MDIC Comex Stat | **Done** | `make mdic-comex-mvp`, `make ci-mdic-comex-mvp` |
| 22 | ANTT pedágios (logística) | **Done** | `make dnit-antt-logistica-mvp`, `make ci-dnit-antt-logistica-mvp` |
| 23 | MAPA ZARC tábua de risco | **Done** | `make mapa-dados-mvp`, `make ci-mapa-dados-mvp` |
| 24 | B3 futuros agro (SOY/CCM/BGI) | **Done** | `make b3-futuros-mvp`, `make ci-b3-futuros-mvp` |

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

**Phase 20 scaffold:** [.local/phases/20-analytics-crossing/README.md](../.local/phases/20-analytics-crossing/README.md) — analytics crossing (feature joins); implementation **not started**.

Remaining tracks (see `.local/PENDING-TASKS.md`):

- Optional: GE bronze referential `cod_ibge` (deferred — post-dbt script is gate)

**Phase 2–3 storage:** `make ci-minio` · `make ci-validate-r2-env` · `make ci-delta-versioning` (native Delta silver versioning)

---

## Related

- [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md) — per-dataset catalog
- [NEW-PROJECT-CHECKLIST.md](NEW-PROJECT-CHECKLIST.md)
- `.local/IMPLEMENTATION-PLAN.md` (gitignored)
- `.local/PENDING-TASKS.md` (gitignored)
