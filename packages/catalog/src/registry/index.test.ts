import { describe, expect, it } from 'vitest';

import { getDatasetById, loadRegistry } from '../index.js';

describe('dataset registry', () => {
  it('loads CONAB MVP datasets', () => {
    const registry = loadRegistry();
    expect(registry.length).toBeGreaterThan(0);
    const estimativa = getDatasetById('conab.estimativa-graos');
    expect(estimativa?.portalLabel).toBe('Estimativa Grãos');
    expect(estimativa?.conabSection).toBe('Produção Agrícola');
    expect(estimativa?.format).toBe('txt');
    expect(estimativa?.schedule).toBe('manual');
  });
});
