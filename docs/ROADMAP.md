# Roadmap — Open Data Agro

> Phases align with `.local/IMPLEMENTATION-PLAN.md` (detailed tasks in `.local/phases/`).

---

## Objective

**Public agro data ingestor** — CONAB portal first, extensible to MAPA, IBGE, INMET.

**Initial source:** [CONAB — Downloads de Arquivos](https://portaldeinformacoes.conab.gov.br/download-arquivos.html)

---

## Stack

| Layer | Technology |
|-------|------------|
| Ingestão | App personalizado (`apps/ingestor`) |
| Data Lake | Cloudflare R2 — Parquet (bronze) |
| Lakehouse | Delta Lake (silver/gold) |
| Processamento | DuckDB |
| Transformação | dbt |
| Qualidade | Great Expectations |
| DB operacional | PostgreSQL |
| DB analítico | DuckDB |

---

## Phase 0 — Scaffold (current)

- [x] Agent harness + governance docs
- [x] `.local` implementation plan + phases
- [ ] Monorepo toolchain (TypeScript, pnpm, Docker Postgres)
- [ ] CI pipeline

## Phase 1–8 — Platform

| Phase | Task | Status |
|-------|------|--------|
| 1 | Ingestor app | Planned |
| 2 | R2 + Parquet bronze | Planned |
| 3 | Delta Lake | Planned |
| 4 | DuckDB processing | Planned |
| 5 | dbt transforms | Planned |
| 6 | Great Expectations | Planned |
| 7 | PostgreSQL ops DB | Planned |
| 8 | DuckDB analytics | Planned |

## Phase 10–14 — CONAB datasets

| Phase | CONAB section | MVP |
|-------|---------------|-----|
| 10 | Produção Agrícola | **P0:** Estimativa Grãos, Série Histórica Grãos |
| 11 | Mercado | **Done** — oferta-demanda, preços (4), mínimos, prohort; `make conab-mercado-full-mvp` |
| 12 | Abastecimento | Planned |
| 13 | Armazenamento e Logística | Planned |
| 14 | Agricultura Familiar | Planned |

---

## MVP path

```
00-platform-scaffold → 01-ingestor-app → 02-data-lake-r2-parquet → 07-db-postgresql → 10-conab-producao-agricola (estimativa-graos)
```

Then: phases 3–6, 8, and remaining CONAB families.

---

## Related

- [.local/IMPLEMENTATION-PLAN.md](../.local/IMPLEMENTATION-PLAN.md) (local, gitignored)
- [NEW-PROJECT-CHECKLIST.md](NEW-PROJECT-CHECKLIST.md)
- [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md)
