/** Branded dataset identifier (e.g. `conab.estimativa-graos`). */
export type DatasetId = string & { readonly __brand: 'DatasetId' };

const DATASET_ID_PATTERN = /^[a-z][a-z0-9]*(\.[a-z][a-z0-9-]*)+$/;

export function isDatasetId(value: string): value is DatasetId {
  return DATASET_ID_PATTERN.test(value);
}

export function parseDatasetId(value: string): DatasetId {
  if (!isDatasetId(value)) {
    throw new Error(`Invalid dataset ID: ${value}`);
  }
  return value;
}
