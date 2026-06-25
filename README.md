# Open Data Agro

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**100% open-source** (MIT) toolkit for Brazilian **agricultural open data** — curated datasets, fetch pipelines, and developer-friendly APIs aligned to official primary sources (MAPA, IBGE, CONAB, INMET, Embrapa, and related `.gov.br` agencies).

> **Status:** documentation and agent harness scaffold — implementation starting.

---

## Mission

Make Brazilian agricultural public data **easy to discover, embed offline, and integrate** — with schemas and transformations traced to official sources, not scraped blogs or unofficial mirrors.

---

## What we build

| Capability | Description |
|------------|-------------|
| **Embed** | Offline snapshots of official agricultural datasets with metadata |
| **Fetch** | Maintainer scripts to refresh data from `.gov.br` APIs and portals |
| **Normalize** | Canonical types, codes, and geospatial joins (IBGE municipalities, CAR, crops) |
| **Expose** | TypeScript library and CLI (planned) for apps, research, and civic tech |

---

## Documentation

| Document | Purpose |
|----------|---------|
| [docs/VISION.md](docs/VISION.md) | Mission, scope, principles |
| [docs/GLOSSARY.md](docs/GLOSSARY.md) | Ubiquitous language |
| [docs/OFFICIAL-SOURCES.md](docs/OFFICIAL-SOURCES.md) | Primary sources per dataset |
| [docs/ROADMAP.md](docs/ROADMAP.md) | Phases and priorities |
| [docs/README.md](docs/README.md) | Full documentation index |

---

## Agent harness

This repo uses the same LLM agent harness as [br-validators](https://github.com/AlexandreZanata/br-validators):

- `agent-rules/` — coding and security rules for agents
- `agent-harness/` — resolve, install, and task-scoped rules
- `.cursor/rules/` — Cursor-specific rule files
- [AGENTS.md](AGENTS.md) — entry point for coding agents

```bash
pip install -r agent-harness/requirements.txt
./agent-harness/resolve-rules.sh data fetch api
```

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Security issues: [SECURITY.md](SECURITY.md) (private advisory).

---

## License

[MIT](LICENSE) — permanently open source. See [docs/OPEN-SOURCE.md](docs/OPEN-SOURCE.md).
