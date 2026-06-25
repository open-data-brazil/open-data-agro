import { describe, expect, it } from 'vitest';
import { isDatasetId, parseDatasetId } from './dataset-id.js';

describe('dataset-id', () => {
  it('accepts valid dataset IDs', () => {
    expect(isDatasetId('conab.estimativa-graos')).toBe(true);
    expect(parseDatasetId('conab.estimativa-graos')).toBe('conab.estimativa-graos');
  });

  it('rejects invalid dataset IDs', () => {
    expect(isDatasetId('invalid')).toBe(false);
    expect(() => parseDatasetId('invalid')).toThrow('Invalid dataset ID');
  });
});
