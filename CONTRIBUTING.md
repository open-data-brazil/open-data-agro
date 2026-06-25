# Contributing to Open Data Agro

Thank you for helping build a **100% open-source** toolkit for Brazilian agricultural open data.

**Language:** All code, comments, docs, commits, and PR descriptions must be in **English**.

---

## Before you start

1. Read [docs/VISION.md](docs/VISION.md) and [docs/GLOSSARY.md](docs/GLOSSARY.md)
2. Check [docs/ROADMAP.md](docs/ROADMAP.md) for planned scope
3. Datasets require an [official source](docs/OFFICIAL-SOURCES.md) — no schema without citation
4. Read [AGENTS.md](AGENTS.md) if using AI coding agents

---

## Open source commitment

This project is **permanently open source** under the [MIT License](LICENSE).

| Allowed | Not allowed |
|---------|-------------|
| MIT-compatible contributions | Proprietary or source-unavailable code |
| Forks and commercial use | "Open core" paid-only datasets in this repo |
| Optional separate UI/adapters packages (also OSS) | CLA that assigns copyright away from contributors |

Full policy: [docs/OPEN-SOURCE.md](docs/OPEN-SOURCE.md)

By contributing, you agree your contributions are licensed under MIT.

---

## How to contribute

### 1. Find or open an issue

- Bug reports: include input, expected vs actual, version
- Features: check roadmap first; new datasets need official sources
- **Security:** see [SECURITY.md](SECURITY.md) — no public issues

### 2. Fork and branch

```bash
git checkout -b feat/conab-harvest-series
# or: fix/ibge-municipality-lookup
```

Branch naming: `feat/`, `fix/`, `docs/`, `test/`, `chore/`

### 3. Implement

- Smallest diff that solves one logical change
- Add/update tests with golden vectors from official docs
- Update [CHANGELOG.md](CHANGELOG.md) under `[Unreleased]`
- Update [docs/OFFICIAL-SOURCES.md](docs/OFFICIAL-SOURCES.md) if adding a dataset
- Follow [docs/API-CONTRACT.md](docs/API-CONTRACT.md) for public exports

### 4. Verify

```bash
make test
make lint
make build
./bin/ingestor --help
```

### 5. Pull request

PR title: `[conab] add harvest forecast embed` or `[docs] add versioning policy`

PR description must include:

- [ ] What changed and why
- [ ] Official source link (if data or schema change)
- [ ] Test vectors added
- [ ] CHANGELOG updated
- [ ] No breaking change OR marked as breaking with migration note
- [ ] MIT-compatible only

---

## Security contributions

Security fixes are welcome and prioritized.

1. **Report first** via [SECURITY.md](SECURITY.md) (private advisory)
2. Wait for maintainer ack before public PR (or coordinate on advisory thread)
3. Include regression test with minimal reproducer
4. Do **not** add real PII or farm production secrets to tests

Security researchers: see severity table in [SECURITY.md](SECURITY.md).

---

## Dataset contributions

Required checklist for new/changed datasets:

- [ ] Row in [docs/OFFICIAL-SOURCES.md](docs/OFFICIAL-SOURCES.md)
- [ ] Business rules in [docs/DATA-RULES.md](docs/DATA-RULES.md) (when created)
- [ ] Use case in `docs/use-cases/` (if new dataset type)
- [ ] Golden vectors in `tests/vectors/`
- [ ] Public API in [docs/API-CONTRACT.md](docs/API-CONTRACT.md)
- [ ] Version impact noted per [docs/VERSIONING.md](docs/VERSIONING.md)

**Wrong agricultural statistics or geospatial joins are data-integrity bugs** — treat with same urgency as [SECURITY.md](SECURITY.md).

---

## Documentation contributions

Docs live in `docs/` and root (`README.md`, `SECURITY.md`, etc.).

- Use ubiquitous terms from [docs/GLOSSARY.md](docs/GLOSSARY.md)
- Link official PDFs/URLs, not secondary blogs
- Keep [docs/README.md](docs/README.md) index updated

---

## Versioning and releases

Maintainers follow [docs/VERSIONING.md](docs/VERSIONING.md).

Contributors:

- **Patch:** bug fix, docs fix, official schema alignment
- **Minor:** new dataset module, new optional export
- **Major:** breaking public API or intentional behavior change on previously valid input

Pre-1.0 (`0.x`): minor bumps may include breaking changes — note in CHANGELOG.

---

## Code of conduct

Be respectful and constructive. See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

---

## Questions

Open a GitHub Discussion (when enabled) or an issue labeled `question`.
