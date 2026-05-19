// Pure-runtime four-state classifier. Expected set is the union of dimType
// bundles for the protocol's registered adapters, plus tvl.

import presets from '../../internal/registry/presets.json' with { type: 'json' }

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
