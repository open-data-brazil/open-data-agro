# Golden vectors — truncated official CONAB downloads (portal headers preserved).

Committed samples are the first 5 data rows plus header from:

- `LevantamentoGraos.txt` → `conab.estimativa-graos`
- `SerieHistoricaGraos.txt` → `conab.serie-historica-graos`

Refresh full files locally (gitignored):

```bash
make conab-reference
```
