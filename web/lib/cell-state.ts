// Pure-runtime four-state classifier. Lives separately from lib/snapshot.ts so
// client components can import it without pulling in node:fs / node:path / node:url.

import { CATEGORIES_EXPECTED } from './categories'

export type CellState = 'na' | 'missing' | 'present' | 'over'

const EXPECTATIONS: Record<string, Record<string, true>> = (() => {
  const out: Record<string, Record<string, true>> = {}
  for (const [category, metrics] of Object.entries(CATEGORIES_EXPECTED)) {
    const set: Record<string, true> = {}
    for (const m of metrics) set[m] = true
    out[category] = set
  }
  return out
})()

// Unseeded categories fall through: present is CellPresent, absent is CellNA.
export function classifyCell(category: string, metric: string, present: boolean): CellState {
  const seed = EXPECTATIONS[category]
  if (!seed) return present ? 'present' : 'na'
  const expected = seed[metric] === true
  if (present && expected) return 'present'
  if (present && !expected) return 'over'
  if (!present && expected) return 'missing'
  return 'na'
}
