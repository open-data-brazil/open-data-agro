# API Contract — Open Data Agro

> **Status:** sketch — refine when `packages/core` is scaffolded.

The public contract will be defined in:

- **This document** — functions, types, errors, versioning
- **[GLOSSARY.md](GLOSSARY.md)** — ubiquitous language
- **[OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md)** — data provenance per module

---

## Design principles

1. **Pure lookup functions** — no network I/O in `@open-data-agro/core`
2. **Explicit results** — `ok: true | false` with typed errors (no thrown exceptions for expected failures)
3. **Tree-shakeable exports** — subpath imports per dataset (`@open-data-agro/core/ibge`)
4. **Metadata accessors** — every dataset exposes `get*Metadata(): DatasetMetadata`

---

## Planned result shape

```typescript
type LookupResult<T> =
  | { ok: true; value: T }
  | { ok: false; error: LookupError };

type DatasetMetadata = {
  datasetId: string;
  capturadoEm: string;       // ISO 8601
  fonteOficial: string;    // URL
  versaoFonte?: string;
  totalRegistros: number;
};
```

---

## Planned modules (Phase 1+)

| Module | Functions (draft) |
|--------|-------------------|
| `ibge` | `getMunicipioPorCodigo`, `searchMunicipios`, `getIbgeMetadata` |
| `conab` | `getSafraPorAno`, `getConabMetadata` |

---

## HTTP adapter (future)

If an HTTP wrapper is added later, create `docs/HTTP-API.md` as a separate adapter layer — core remains offline embeds only.

---

## Versioning

Public API follows [VERSIONING.md](VERSIONING.md). Breaking changes require MAJOR bump after `1.0.0`.
