// Pure-runtime four-state classifier. Lives separately from lib/snapshot.ts so
// client components can import it without pulling in node:fs / node:path / node:url.

import presets from '../../internal/registry/presets.json' with { type: 'json' }

export type CellState = 'na' | 'missing' | 'present' | 'over'

// Spotlight excludes 'na' because masking non-matching cells to na would leave
// an all-na matrix indistinguishable from the spotlight target.
export const SPOTLIGHT_STATES = ['present', 'missing', 'over'] as const
export type SpotlightState = (typeof SPOTLIGHT_STATES)[number]

export function parseSpotlightParam(raw: string | null | undefined): SpotlightState | '' {
  const v = raw ?? ''
  return (SPOTLIGHT_STATES as readonly string[]).includes(v) ? (v as SpotlightState) : ''
}

// Same JSON source as Go and tools/build-snapshot.ts.
const EXPECTATIONS: Record<string, Record<string, true>> = (() => {
  const out: Record<string, Record<string, true>> = {}
  for (const [category, metrics] of Object.entries(presets.categories)) {
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
