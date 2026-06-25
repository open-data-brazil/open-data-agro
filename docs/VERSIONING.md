# Versioning policy

> **Open Data Agro** follows [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html) (SemVer).
> Version applies to npm releases (when packages exist) and git tags.

---

## Version format

```
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
```

| Segment | When to increment | Example |
|---------|-------------------|---------|
| **MAJOR** | Breaking public API or intentional behavior change | `1.0.0` → `2.0.0` |
| **MINOR** | New backward-compatible functionality | `0.2.0` — new `conab` module |
| **PATCH** | Backward-compatible bug fix | `0.1.1` — fix IBGE code lookup |
| **PRERELEASE** | Pre-release quality | `0.1.0-alpha.1` |

---

## Pre-1.0 policy (`0.x.y`)

During alpha/beta:

- **`0.MINOR.PATCH`** — public API may change between minors
- Breaking changes documented in [CHANGELOG.md](../CHANGELOG.md)
- Consumers should pin exact version: `"@open-data-agro/core": "0.1.0"`

Planned milestones:

| Version | Meaning |
|---------|---------|
| `0.1.0-alpha` | IBGE municipalities + harness + docs |
| `0.2.0` | CONAB safras embed |
| `0.3.0` | MAPA catalog (initial datasets) |
| `1.0.0` | **Stable API contract** — SemVer guarantees apply fully |

---

## Data versioning

Embedded datasets version with the package release that contains them.

- **PATCH** — correction aligned with official source, same schema
- **MINOR** — new rows or optional fields, backward compatible
- **MAJOR** — breaking schema or code remapping (document migration)

Always record `capturadoEm` in dataset metadata.

---

## Support window

| Version | Supported |
|---------|-----------|
| Latest `1.x` | Yes |
| Previous minor | Security and data-integrity fixes only |
| `0.x` | Best effort |
| Unreleased `main` | Fix forward |

---

## Related

- [CHANGELOG.md](../CHANGELOG.md)
- [API-CONTRACT.md](API-CONTRACT.md)
- [SECURITY.md](../SECURITY.md)
