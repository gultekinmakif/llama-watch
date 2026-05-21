// Expected set is tvl + CATEGORIES_EXPECTED[category]. dimType bundles stay
// available for column narrowing on the adapter dropdown.

import presets from '../../internal/registry/presets.json' with { type: 'json' }
import { CATEGORIES_EXPECTED } from './categories'

export type CellState = 'na' | 'missing' | 'present' | 'unexpected'

export function metricsForDimType(dimType: string): readonly string[] {
  const metrics = (presets as Record<string, readonly string[]>)[dimType]
  return metrics ?? []
}

export const DIMTYPE_KEYS: readonly string[] = Object.keys(presets)

export function classifyCell(
  category: string | undefined,
  metric: string,
  present: boolean,
): CellState {
  const expectedSet = category != null ? CATEGORIES_EXPECTED[category] : undefined
  const expected =
    metric === 'tvl' || (expectedSet != null && expectedSet.includes(metric as never))
  if (present && expected) return 'present'
  if (present && !expected) return 'unexpected'
  if (!present && expected) return 'missing'
  return 'na'
}
