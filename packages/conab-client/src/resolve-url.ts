import { getDatasetById, type DatasetRegistryEntry } from '@open-data-agro/catalog';

const CONAB_PORTAL_BASE = 'https://portaldeinformacoes.conab.gov.br';

export type ResolvedDownloadUrl = {
  url: string;
  portalPage: string;
  discoveredAt: string;
};

export function resolveDownloadUrl(datasetId: string): ResolvedDownloadUrl {
  const entry = getDatasetById(datasetId);
  if (!entry) {
    throw new Error(`No CONAB download mapping for dataset: ${datasetId}`);
  }
  return {
    url: entry.sourceUrl,
    portalPage: `${CONAB_PORTAL_BASE}/download-arquivos.html`,
    discoveredAt: entry.discoveredAt,
  };
}

export function getRegistryEntry(datasetId: string): DatasetRegistryEntry {
  const entry = getDatasetById(datasetId);
  if (!entry) {
    throw new Error(`No CONAB registry entry for dataset: ${datasetId}`);
  }
  return entry;
}
