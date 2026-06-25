# Golden vectors — truncated official CONAB downloads (portal headers preserved).

Committed samples are the first 5 data rows plus header from:

- `LevantamentoGraos.txt` → `conab.estimativa-graos`
- `SerieHistoricaGraos.txt` → `conab.serie-historica-graos`
- `OfertaDemanda.txt` → `conab.oferta-demanda`
- `Estoques.txt` → `conab.estoques-publicos`

Refresh full files locally (gitignored):

```bash
make conab-reference
```
