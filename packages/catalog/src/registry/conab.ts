import { parseDatasetId } from '../dataset-id.js';
import type { DatasetRegistryEntry } from '../dataset-registry-entry.js';

const PORTAL_BASE = 'https://portaldeinformacoes.conab.gov.br';
const PORTAL_PAGE = `${PORTAL_BASE}/download-arquivos.html`;
const DISCOVERED_AT = '2026-06-25T00:00:00.000Z';

function conabEntry(
  datasetId: string,
  portalFile: string,
  portalLabel: string,
  conabSection: string,
  format: DatasetRegistryEntry['format'] = 'txt',
  delimiter = ';',
): DatasetRegistryEntry {
  return {
    datasetId: parseDatasetId(datasetId),
    sourceUrl: `${PORTAL_BASE}/downloads/arquivos/${portalFile}`,
    format,
    schedule: 'manual',
    conabSection,
    portalLabel,
    delimiter,
    discoveredAt: DISCOVERED_AT,
  };
}

/** CONAB datasets resolved from download-arquivos.html (2026-06-25). */
export const CONAB_REGISTRY: DatasetRegistryEntry[] = [
  conabEntry('conab.custo-producao', 'CustoProducao.txt', 'Custo de Produção', 'Produção Agrícola'),
  conabEntry('conab.estimativa-cana', 'LevantamentoCana.txt', 'Estimativa Cana-de-Açúcar', 'Produção Agrícola'),
  conabEntry('conab.estimativa-graos', 'LevantamentoGraos.txt', 'Estimativa Grãos', 'Produção Agrícola'),
  conabEntry('conab.serie-historica-cana', 'SerieHistoricaCana.txt', 'Serie Histórica Cana-de-Açúcar', 'Produção Agrícola'),
  conabEntry('conab.serie-historica-graos', 'SerieHistoricaGraos.txt', 'Serie Histórica Grãos', 'Produção Agrícola'),
  conabEntry('conab.estimativa-cafe', 'SerieHistoricaCafe.txt', 'Estimativa Café', 'Produção Agrícola'),
  conabEntry('conab.serie-historica-cafe', 'LevantamentoCafe.txt', 'Serie Histórica Café', 'Produção Agrícola'),
  conabEntry('conab.armazenagem', 'ArmazensCadastrados.txt', 'Armazenagem', 'Armazenamento e Logística'),
  conabEntry('conab.frete', 'Frete.txt', 'Frete', 'Armazenamento e Logística'),
  conabEntry('conab.serie-historica-capacidade-estatica', 'exportacao_capacidade_estatica.xls', 'Série Histórica da Capacidade Estática', 'Armazenamento e Logística', 'xlsx'),
  conabEntry('conab.estoques-publicos', 'Estoques.txt', 'Estoques Públicos', 'Abastecimento'),
  conabEntry('conab.operacoes-comercializacao', 'Leilao.txt', 'Operações de Comercialização', 'Abastecimento'),
  conabEntry('conab.vendas-balcao', 'VendaBalcao.txt', 'Vendas em Balcão', 'Abastecimento'),
  conabEntry('conab.oferta-demanda', 'OfertaDemanda.txt', 'Oferta e Demanda', 'Mercado'),
  conabEntry('conab.precos-minimos', 'PrecoMinimo.txt', 'Preços Mínimos', 'Mercado'),
  conabEntry('conab.precos-agropecuarios-mensal-uf', 'PrecosMensalUF.txt', 'Preços agropecuários Mensal UF', 'Mercado'),
  conabEntry('conab.precos-agropecuarios-mensal-municipio', 'PrecosMensalMunicipio.txt', 'Preços agropecuários Mensal Município', 'Mercado'),
  conabEntry('conab.precos-agropecuarios-semanal-uf', 'PrecosSemanalUF.txt', 'Preços agropecuários Semanal UF', 'Mercado'),
  conabEntry('conab.precos-agropecuarios-semanal-municipio', 'PrecosSemanalMunicipio.txt', 'Preços agropecuários Semanal Municipio', 'Mercado'),
  conabEntry('conab.prohort-diario', 'ProhortDiario.txt', 'Prohort Diário', 'Mercado'),
  conabEntry('conab.prohort-mensal', 'ProhortMensal.txt', 'Prohort Mensal', 'Mercado'),
  conabEntry('conab.alimenta-brasil-entregas', 'PAA_Entregas.txt', 'Programa Alimenta Brasil - Entregas', 'Agricultura Familiar'),
  conabEntry('conab.alimenta-brasil-propostas', 'PAA_PropostaFormalizadasExecutada.txt', 'Programa Alimenta Brasil - Propostas', 'Agricultura Familiar'),
];

export const CONAB_PORTAL_PAGE = PORTAL_PAGE;
