import { describe, expect, it } from 'vitest';

import { resolveDownloadUrl } from './resolve-url.js';

describe('resolveDownloadUrl', () => {
  it('resolves estimativa-graos from catalog', () => {
    const resolved = resolveDownloadUrl('conab.estimativa-graos');
    expect(resolved.url).toBe(
      'https://portaldeinformacoes.conab.gov.br/downloads/arquivos/LevantamentoGraos.txt',
    );
    expect(resolved.portalPage).toContain('download-arquivos.html');
  });

  it('throws for unknown dataset', () => {
    expect(() => resolveDownloadUrl('unknown.dataset')).toThrow('No CONAB download mapping');
  });
});
