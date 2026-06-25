# New Project Checklist

> Complete **before writing the first line of code**.
> Status: documentation and harness phase — ready for toolchain scaffold.

---

## Architecture and domain

- [x] **Layers defined** — fetch / embed / core / adapters — see [ARCHITECTURE.md](ARCHITECTURE.md)
- [x] **Value types** — canonical codes and metadata — see [GLOSSARY.md](GLOSSARY.md)
- [ ] **Business rules** — GIVEN/WHEN/THEN in `DATA-RULES.md` (create with first dataset)
- [ ] **Use cases** — expand beyond UC-001 in [use-cases/](use-cases/)
- [x] **API contract** — [API-CONTRACT.md](API-CONTRACT.md)
- [x] **Glossary** — [GLOSSARY.md](GLOSSARY.md)

---

## Security (OWASP)

- [x] **SECURITY.md** — private vulnerability reporting
- [x] **SECURITY-PRACTICES.md** — maintainer and integrator guidance
- [ ] **OWASP Top 10:2025** — dependency audit when scaffold exists
- [x] **Agentic 2026** — harness rules apply to agent sessions

---

## Governance and releases

- [x] **LICENSE** — MIT
- [x] **OPEN-SOURCE.md** — 100% open source commitment
- [x] **VERSIONING.md** — SemVer and release process
- [x] **CONTRIBUTING.md** — contribution + security contribution
- [x] **CODE_OF_CONDUCT.md**
- [x] **CHANGELOG.md** — Keep a Changelog format
- [x] **GOVERNANCE.md** — document index

---

## Official sources

- [x] **Source catalog scaffold** — [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md)
- [ ] **Per-dataset rows** — fill as implemented

---

## Agent harness

- [x] **Harness installed** — `agent-rules/`, `agent-harness/`, `.cursor/rules/`
- [x] **AGENTS.md** — project entry point

---

## Implementation readiness

- [ ] Toolchain chosen — **TypeScript 5+**, pnpm, Vitest (see [ARCHITECTURE.md](ARCHITECTURE.md))
- [ ] Monorepo scaffold (`packages/core`, `scripts/`)
- [ ] Golden test vectors from official API samples
- [ ] CI pipeline

---

## Sign-off

| Role | Name | Date |
|------|------|------|
| Product / domain | | |
| Tech lead | | |

When implementation checklist items are checked, coding may begin.
