import { loadRegistry } from '@open-data-agro/catalog';

export function listCatalogEntries() {
  return loadRegistry();
}
