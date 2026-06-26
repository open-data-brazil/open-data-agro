# New Project Checklist

> Pre-coding checklist — **updated 2026-06-26** after collection sprint (Phases 0–19).  
> Status: **Go local-first platform operational** — 47 catalog datasets with offline `make *-mvp` and CI gates.

---

## Architecture and domain

- [x] **Layers defined** — ingest / bronze / silver / gold / analytics — see [ARCHITECTURE.md](ARCHITECTURE.md)
- [x] **Value types** — canonical codes and metadata — see [GLOSSARY.md](GLOSSARY.md)
- [—] **Business rules** — `DATA-RULES.md` deferred to Phase 20+ analytics (cross-mart joins); bronze GE + `validate_codigo_ibge.py` cover collection gates
- [~] **Use cases** — [UC-001](use-cases/UC-001-lookup-municipality.md) exists; expand when public API ships
- [x] **API contract** — [API-CONTRACT.md](API-CONTRACT.md) (CLI / lake layout sketch)
- [x] **Glossary** — [GLOSSARY.md](GLOSSARY.md)

---

## Security (OWASP)

- [x] **SECURITY.md** — private vulnerability reporting
- [x] **SECURITY-PRACTICES.md** — maintainer and integrator guidance
- [x] **OWASP Top 10:2025** — `govulncheck` in GitHub Actions `security` job
- [x] **Agentic 2026** — harness rules apply to agent sessions

---

## Governance and releases

- [x] **LICENSE** — MIT
- [x] **OPEN-SOURCE.md** — 100% open source commitment
- [x] **VERSIONING.md** — SemVer and release process
- [x] **CONTRIBUTING.md** — contribution + security contribution
- [x] **CODE_OF_CONDUCT.md**
- [x] **CHANGELOG.md** — Keep a Changelog format (`[Unreleased]` tracks collection sprint)
- [x] **GOVERNANCE.md** — document index

---

## Official sources

- [x] **Source catalog** — [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md) (47 datasets, `**Pn — implemented**` status)
- [x] **Per-dataset rows** — synced with `configs/catalog/`; gate: `python3 scripts/ci/check_official_sources_status.py`

---

## Agent harness

- [x] **Harness installed** — `agent-rules/`, `agent-harness/`, `.cursor/rules/`
- [x] **AGENTS.md** — project entry point

---

## Implementation readiness

- [x] **Toolchain** — **Go 1.22+**, DuckDB, dbt, Python quality scripts (see [ROADMAP.md](ROADMAP.md) stack table)
- [x] **Repo layout** — `cmd/ingestor`, `cmd/processor`, `internal/`, `dbt/`, `expectations/`, `lake/` local-first
- [x] **Golden test vectors** — `internal/*/testdata/` + `go test ./internal/ingest/`
- [x] **CI pipeline** — [`.github/workflows/ci.yml`](../.github/workflows/ci.yml); local mirror: `make ci-go`, `make ci-dbt`

**Collection sprint exit gates:**

```bash
make collection-full-mvp      # local
make ci-collection-full-mvp   # offline CI (also in GitHub Actions dbt job)
```

---

## Sign-off

| Role | Name | Date |
|------|------|------|
| Product / domain | Collection sprint (João P1) | 2026-06-26 |
| Tech lead | Phases 0–19 E2E + CI gates | 2026-06-26 |

Post-collection work: see `.local/PENDING-TASKS.md` (Phase 20 analytics, R2/MinIO production track).
