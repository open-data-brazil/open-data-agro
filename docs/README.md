# Documentation index — Open Data Agro

> **100% open-source** toolkit for Brazilian agricultural open data.
> All docs in English. Official sources only — no guesswork on schemas or codes.

---

## Start here

| Document | Purpose |
|----------|---------|
| [VISION.md](VISION.md) | Mission, scope, principles, non-goals |
| [GLOSSARY.md](GLOSSARY.md) | Ubiquitous language — use these terms in code |
| [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md) | Primary sources per dataset (MAPA, IBGE, CONAB, INMET) |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Module layout, layers, design patterns |
| [API-CONTRACT.md](API-CONTRACT.md) | Public API contract (functions, types, errors) |
| [ROADMAP.md](ROADMAP.md) | Phases and priorities |
| [POSTGRES-UNIFIED-SYNC.md](POSTGRES-UNIFIED-SYNC.md) | Gold marts → PostgreSQL `analytics` schema (Stage G) |
| [NEW-PROJECT-CHECKLIST.md](NEW-PROJECT-CHECKLIST.md) | Pre-implementation checklist |

## Governance (license, security, releases)

| Document | Purpose |
|----------|---------|
| [OPEN-SOURCE.md](OPEN-SOURCE.md) | 100% MIT open source commitment |
| [VERSIONING.md](VERSIONING.md) | SemVer, releases, support window |
| [SECURITY-PRACTICES.md](SECURITY-PRACTICES.md) | Maintainer & integrator security |
| [GOVERNANCE.md](GOVERNANCE.md) | Document map and decision log |
| [../LICENSE](../LICENSE) | MIT license text |
| [../SECURITY.md](../SECURITY.md) | Report vulnerabilities (private) |
| [../CONTRIBUTING.md](../CONTRIBUTING.md) | How to contribute |
| [../CHANGELOG.md](../CHANGELOG.md) | Release history |

## Use cases

| ID | File | Summary |
|----|------|---------|
| UC-001 | [use-cases/UC-001-lookup-municipality.md](use-cases/UC-001-lookup-municipality.md) | Resolve IBGE municipality code for ag data joins |

## For agents

1. Read [GLOSSARY.md](GLOSSARY.md) before naming anything in code.
2. Implement only schemas traced to [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md).
3. Check [API-CONTRACT.md](API-CONTRACT.md) before adding public exports.
4. Run `./agent-harness/resolve-rules.sh` for task-specific rules.
