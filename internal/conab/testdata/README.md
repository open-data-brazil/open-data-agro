# Golden vectors — truncated official CONAB downloads (portal headers preserved).

Committed samples are the first 5 data rows plus header from:

- `LevantamentoGraos.txt` → `conab.estimativa-graos`
- `SerieHistoricaGraos.txt` → `conab.serie-historica-graos`
- `OfertaDemanda.txt` → `conab.oferta-demanda`
- `PrecoMinimo.txt` → `conab.precos-minimos`
- `PrecosSemanalUF.txt` → `conab.precos-agropecuarios-semanal-uf`
- `PrecosSemanalMunicipio.txt` → `conab.precos-agropecuarios-semanal-municipio`
- `PrecosMensalUF.txt` → `conab.precos-agropecuarios-mensal-uf`
- `PrecosMensalMunicipio.txt` → `conab.precos-agropecuarios-mensal-municipio`
- `Estoques.txt` → `conab.estoques-publicos`
- `ArmazensCadastrados.txt` → `conab.armazenagem`
- `Frete.txt` → `conab.frete`
- `exportacao_capacidade_estatica.xls` → `conab.serie-historica-capacidade-estatica`
- `PAA_Entregas.txt` → `conab.alimenta-brasil-entregas`
- `PAA_PropostaFormalizadasExecutada.txt` → `conab.alimenta-brasil-propostas`

Refresh full files locally (gitignored):

```bash
make conab-reference
```
