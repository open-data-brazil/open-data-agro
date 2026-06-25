import { isDatasetId, type DatasetId } from '../dataset-id.js';
import type { DatasetRegistryEntry } from '../dataset-registry-entry.js';
import { CONAB_REGISTRY } from './conab.js';

const REGISTRY: DatasetRegistryEntry[] = [...CONAB_REGISTRY];

export function loadRegistry(): DatasetRegistryEntry[] {
  return [...REGISTRY];
}

export function getDatasetById(datasetId: string): DatasetRegistryEntry | undefined {
  if (!isDatasetId(datasetId)) {
    return undefined;
  }
  return REGISTRY.find((entry) => entry.datasetId === datasetId);
}

export function requireDataset(datasetId: string): DatasetRegistryEntry {
  const entry = getDatasetById(datasetId);
  if (!entry) {
    throw new Error(`Unknown dataset: ${datasetId}`);
  }
  return entry;
}

export function listDatasetIds(): DatasetId[] {
  return REGISTRY.map((entry) => entry.datasetId);
}
