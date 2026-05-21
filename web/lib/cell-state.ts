// Two classifiers. The default `classifyCell` reads the protocol's own adapter
// bundles; the matrix ships in this mode so coverage stays at the
// per-protocol-emission bar. `classifyByCategory` is the alt mode the UI flips
// to via `?mode=dimensions`, surfacing protocols missing dimension adapters
// against the category's curated expected set.

import presets from '../../internal/registry/presets.json' with { type: 'json' }
import { CATEGORIES_EXPECTED } from './categories'

export type CellState = 'na' | 'missing' | 'present' | 'unexpected'

const BUNDLES: Record<string, Record<string, true>> = (() => {
  const out: Record<string, Record<string, true>> = {}
  for (const [dt, metrics] of Object.entries(presets)) {
    const set: Record<string, true> = {}
    for (const m of metrics) set[m] = true
    out[dt] = set
  }
  return out
})()

export function metricsForDimType(dimType: string): readonly string[] {
  const metrics = (presets as Record<string, readonly string[]>)[dimType]
  return metrics ?? []
}

export const DIMTYPE_KEYS: readonly string[] = Object.keys(presets)

export function classifyCell(
  dimTypes: readonly string[],
  metric: string,
  present: boolean,
): CellState {
  let expected = metric === 'tvl'
  if (!expected) {
    for (const dt of dimTypes) {
      if (BUNDLES[dt]?.[metric]) {
        expected = true
        break
      }
    }
  }
  if (present && expected) return 'present'
  if (present && !expected) return 'unexpected'
  if (!present && expected) return 'missing'
  return 'na'
}

export function classifyByCategory(
  category: string | undefined,
  metric: string,
  present: boolean,
): CellState {
  const expectedSet = category != null ? CATEGORIES_EXPECTED[category] : undefined
  // skip categories with no expected fields.
  if (expectedSet == null) return present ? 'present' : 'na'
  const expected =
    metric === 'tvl' || (expectedSet as readonly string[]).includes(metric)
  if (present && expected) return 'present'
  if (present && !expected) return 'unexpected'
  if (!present && expected) return 'missing'
  return 'na'
}
