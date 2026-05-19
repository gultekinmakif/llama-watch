// Pure-runtime helpers backing the category and dimType filter dropdowns.
// No node imports so client components can consume these directly.

import presets from '../../internal/registry/presets.json' with { type: 'json' }
import { CATEGORIES_EXPECTED } from './categories'
import type { ColumnKey } from './snapshot'

export const CATEGORIES: readonly string[] = Object.keys(CATEGORIES_EXPECTED).sort()

export const DIMTYPES: readonly string[] = Object.keys(presets).sort()

// Returns a fresh array so callers can mutate without poisoning the source.
export function expectedColumnsFor(category: string): ColumnKey[] {
  const metrics = CATEGORIES_EXPECTED[category]
  if (!metrics) return []
  return [...metrics] as ColumnKey[]
}

// Returns a fresh array so callers can mutate without poisoning the source.
export function metricsFor(dimType: string): ColumnKey[] {
  const metrics = (presets as Record<string, readonly string[]>)[dimType]
  if (!metrics) return []
  return [...metrics] as ColumnKey[]
}
