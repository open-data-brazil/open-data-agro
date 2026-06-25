# Architecture — Open Data Agro

> **TypeScript monorepo (planned):** core library + fetch scripts + optional CLI/docs.
> Offline embeds at runtime; network only in maintainer fetch scripts.

---

## Planned monorepo layout

```
open-data-agro/
├── packages/
│   └── core/                       # npm — MIT, minimal runtime deps
│       ├── src/
│       │   ├── ibge/
│       │   ├── conab/
│       │   ├── types/
│       │   │   ├── lookup-result.ts
│       │   │   └── dataset-metadata.ts
│       │   └── index.ts
│       └── data/                   # embedded JSON (generated, not hand-edited)
├── scripts/
│   ├── fetch-ibge-municipios.ts
│   ├── fetch-conab-safras.ts
│   └── data-refresh-bot.ts       # planned
├── data/
│   └── refresh-reports/            # drift logs (planned)
├── agent-rules/                    # LLM coding rules (harness)
├── agent-harness/
├── docs/
└── pnpm-workspace.yaml
```

---

## Dependency graph

```
apps/cli (planned)     ──► packages/core
apps/docs (planned)    ──► packages/core

packages/core        ──► embedded data only (no runtime network)

scripts/fetch-*      ──► official APIs (maintainer only)
```

Normalization logic exists **once** in `packages/core`. Fetch scripts write embedded artifacts.

---

## Layers

| Layer | Responsibility | Depends on |
|-------|----------------|------------|
| **embed** | Versioned JSON/Parquet snapshots + `metadata.json` | Nothing |
| **core** | Typed lookups, joins, validation of codes | embed |
| **fetch** | Download and transform from official sources | Nothing at runtime |
| **adapters** (future) | Optional live HTTP fallback | core |

No network, filesystem reads of arbitrary paths, or env vars inside `core/` lookup functions.

---

## Data catalog pattern

Each dataset registers:

- `datasetId` — stable string (`ibge.municipios`)
- `metadata.json` — [GLOSSARY.md](GLOSSARY.md) `DatasetMetadata`
- `fetch` script — `scripts/fetch-*.ts`
- `vectors` — `tests/vectors/<dataset>.official.json`

---

## Agent harness

Same pattern as [br-validators](https://github.com/AlexandreZanata/br-validators):

| Path | Purpose |
|------|---------|
| `agent-rules/` | Full LLM best-practices rule tree |
| `agent-harness/` | `resolve-rules.sh`, `rules-path.sh`, install scripts |
| `.cursor/rules/` | Cursor `alwaysApply` rules |
| `AGENTS.md` | Agent session entry point |

---

## Related

- [API-CONTRACT.md](API-CONTRACT.md)
- [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md)
- [ROADMAP.md](ROADMAP.md)
