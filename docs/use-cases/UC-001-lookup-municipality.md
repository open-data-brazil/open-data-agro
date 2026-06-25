# UC-001 — Lookup IBGE municipality

## Summary

Resolve a seven-digit IBGE municipality code to canonical name, UF, and region for agricultural data joins.

## Actors

- Application developer
- Data pipeline

## Preconditions

- `ibge.municipios` dataset embedded in `@open-data-agro/core`
- Input is a string of seven digits

## Main flow

1. **GIVEN** a valid IBGE code `3550308`
2. **WHEN** `getMunicipioPorCodigo('3550308')` is called
3. **THEN** result is `{ ok: true, value: { codigo: '3550308', nome: 'São Paulo', uf: 'SP', ... } }`

## Alternate flows

| Case | Expected |
|------|----------|
| Unknown code | `{ ok: false, error: 'NOT_FOUND' }` |
| Invalid length | `{ ok: false, error: 'INVALID_FORMAT' }` |

## Official source

[IBGE API de localidades](https://servicodados.ibge.gov.br/api/docs/localidades) — see [OFFICIAL-SOURCES.md](../OFFICIAL-SOURCES.md).

## Notes

- Display names may retain accents; codes are ASCII digits only.
- Golden vectors: `tests/vectors/ibge.municipios.official.json` (planned).
