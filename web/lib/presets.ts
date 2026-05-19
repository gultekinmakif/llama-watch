// Pure-runtime helpers backing the category and dimType filter dropdowns.
// No node imports so client components can consume these directly.

import presets from '../../internal/registry/presets.json' with { type: 'json' }
import type { ColumnKey } from './snapshot'

// Same JSON source as Go and tools/build-snapshot.ts.

export const CATEGORIES: readonly string[] = Object.keys(presets.categories).sort()

export const DIMTYPES: readonly string[] = Object.keys(presets.dimTypes).sort()

// Returns a fresh array so callers can mutate without poisoning the source.
export function expectedColumnsFor(category: string): ColumnKey[] {
  const metrics = (presets.categories as Record<string, readonly string[]>)[category]
  if (!metrics) return []
  return [...metrics] as ColumnKey[]
}

// Returns a fresh array so callers can mutate without poisoning the source.
export function metricsFor(dimType: string): ColumnKey[] {
  const metrics = (presets.dimTypes as Record<string, readonly string[]>)[dimType]
  if (!metrics) return []
  return [...metrics] as ColumnKey[]
}
