import type { DatasetId } from './dataset-id.js';

/** Capture metadata for an embedded or ingested dataset. */
export type DatasetMetadata = {
  datasetId: DatasetId;
  capturadoEm: string;
  fonteOficial: string;
  versaoFonte?: string;
  totalRegistros: number;
};
