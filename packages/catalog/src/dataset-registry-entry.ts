import type { DatasetId } from './dataset-id.js';

export type DatasetFormat = 'csv' | 'xlsx' | 'txt';

export type DatasetRegistryEntry = {
  datasetId: DatasetId;
  sourceUrl: string;
  format: DatasetFormat;
  schedule: string;
  conabSection: string;
  portalLabel: string;
  /** CSV/TXT delimiter when format is csv or txt */
  delimiter?: string;
  discoveredAt: string;
};
