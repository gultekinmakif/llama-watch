import presets from '../../internal/registry/presets.json' with { type: 'json' }
import { CATEGORIES_EXPECTED } from './categories'
import type { ColumnKey } from './snapshot'

export type CellState = 'na' | 'missing' | 'present' | 'unexpected'

type Classifier = (category: string | undefined, metric: ColumnKey, present: boolean) => CellState

interface ReclassifyTarget {
  category?: string
  cells: Record<ColumnKey, CellState>
}

export interface ReclassifyResult {
  cells: Record<ColumnKey, CellState>
  coverage: number
  expected: number
}

export function reclassifyRow(
  row: ReclassifyTarget,
  columns: readonly { key: ColumnKey }[],
  classify: Classifier,
): ReclassifyResult {
  const cells = {} as Record<ColumnKey, CellState>
  let coverage = 0
  let expected = 0
  for (const col of columns) {
    const was = row.cells[col.key]
    const isPresent = was === 'present' || was === 'unexpected'
    const state = classify(row.category, col.key, isPresent)
    cells[col.key] = state
    if (state === 'present') coverage += 1
    if (state === 'present' || state === 'missing') expected += 1
  }
  return { cells, coverage, expected }
}

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
  let expected = false
  for (const dt of dimTypes) {
    if (BUNDLES[dt]?.[metric]) {
      expected = true
      break
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
  const expected = (expectedSet as readonly string[]).includes(metric)
  if (present && expected) return 'present'
  if (present && !expected) return 'unexpected'
  if (!present && expected) return 'missing'
  return 'na'
}
