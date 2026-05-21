// Build-time snapshot reader. Projects the wire shape {cells, protocols}
// onto the closed column set the table renders. Keep COLUMNS in lockstep
// with internal/registry/columns.go.

import { readFileSync, statSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

// Resolve relative to this source file so the path holds no matter where the
// Next build is invoked from (CI uses /, local dev uses /web, both work).
const HERE = dirname(fileURLToPath(import.meta.url))
const RESOLVED_SNAPSHOT_PATH = resolve(HERE, '..', '..', 'var', 'snapshot', 'snapshot.json')

const COLUMNS = [
  { key: 'tvl', label: 'TVL' },
  { key: 'dailyFees', label: 'Daily Fees' },
  { key: 'dailyRevenue', label: 'Daily Revenue' },
  { key: 'dailyVolume', label: 'Daily Volume' },
  { key: 'dailyNotionalVolume', label: 'Notional Volume' },
  { key: 'dailyPremiumVolume', label: 'Premium Volume' },
  { key: 'openInterestAtEnd', label: 'Open Interest' },
  { key: 'dailyBridgeVolume', label: 'Bridge Volume' },
  { key: 'dailyActiveUsers', label: 'Active Users' },
  { key: 'dailyUserFees', label: 'User Fees' },
  { key: 'dailyHoldersRevenue', label: 'Holders Revenue' },
  { key: 'dailyProtocolRevenue', label: 'Protocol Revenue' },
  { key: 'dailySupplySideRevenue', label: 'Supply-Side Revenue' },
  { key: 'dailyCreatorRevenue', label: 'Creator Revenue' },
  { key: 'dailyBribesRevenue', label: 'Bribes Revenue' },
  { key: 'dailyTokenTaxes', label: 'Token Taxes' },
  { key: 'longOpenInterestAtEnd', label: 'Long Open Interest' },
  { key: 'shortOpenInterestAtEnd', label: 'Short Open Interest' },
  { key: 'dailyTransactionsCount', label: 'Transactions' },
  { key: 'dailyGasUsed', label: 'Gas Used' },
  { key: 'dailyNewUsers', label: 'New Users' },
  { key: 'dailyNormalizedVolume', label: 'Normalized Volume' },
  { key: 'dailyActiveLiquidity', label: 'Active Liquidity' },
  { key: 'tokenIncentives', label: 'Token Incentives' },
] as const

export type ColumnKey = (typeof COLUMNS)[number]['key']

export interface Column {
  key: ColumnKey
  label: string
}

import { classifyCell, classifyByCategory, reclassifyRow, type CellState } from './cell-state'

export type Cells = Record<ColumnKey, CellState>

export interface Row {
  slug: string
  name: string
  category?: string
  chains: string[]
  cells: Cells
  // The dimType adapters the protocol is registered with, post-manifest-filter.
  // Passed to the classifier; also surfaced so the static detail page can render four states.
  dimTypes: string[]
  // Precomputed at build time so sort/render does not redo Object.values on every pass.
  coverage: number
  // present + missing for this row: the size of the protocol's expected-metric set.
  expected: number
}

export interface SnapshotStats {
  tracked: number
  coveragePct: number
  // Re-scored against CATEGORIES_EXPECTED; HeroStrip picks this on ?mode=dimensions.
  coveragePctDimensions: number
  updatedAt: string
}

export interface Snapshot {
  columns: Column[]
  rows: Row[]
  total: number
  stats: SnapshotStats
}

export interface RawCell {
  slug: string
  metric: string
  codePath: string
}

export interface RawProtocol {
  slug: string
  name: string
  category?: string
  chains: string[]
  dataFile: string
  dimTypes: string[]
}

export interface RawSnapshot {
  cells: RawCell[]
  protocols: RawProtocol[]
  generatedAt?: string // by tools/build-snapshot.ts
}

// Runs at build time only; sync fs is correct here.
export function loadSnapshot(): Snapshot {
  const raw = readRaw()
  const presence = presenceBySlug(raw.cells)
  const rows = raw.protocols.map((p) => projectRow(p, presence.get(p.slug)))
  if (rows.length === 0) {
    // Next refuses output:'export' on dynamic routes when generateStaticParams
    // returns []; catching here keeps the failure mode loud and on-topic.
    throw new Error(
      `snapshot at ${RESOLVED_SNAPSHOT_PATH} has zero protocols; static export needs at least one`,
    )
  }
  // Densest columns first so the matrix opens with the highest-coverage metrics on the left.
  const presentCounts = new Map<ColumnKey, number>()
  for (const col of COLUMNS) presentCounts.set(col.key, 0)
  for (const r of rows) {
    for (const col of COLUMNS) {
      if (r.cells[col.key] === 'present') {
        presentCounts.set(col.key, (presentCounts.get(col.key) ?? 0) + 1)
      }
    }
  }
  const columns = COLUMNS.map((c) => ({ ...c })).sort(
    (a, b) => (presentCounts.get(b.key) ?? 0) - (presentCounts.get(a.key) ?? 0),
  )

  // Coverage = present cells / (present + missing) cells across the matrix.
  // 'na' and 'unexpected' are excluded so the percentage reflects how much of the
  // expected data we actually see.
  let presentTotal = 0
  let trackedTotal = 0
  for (const r of rows) {
    for (const col of COLUMNS) {
      const s = r.cells[col.key]
      if (s === 'present') {
        presentTotal += 1
        trackedTotal += 1
      } else if (s === 'missing') {
        trackedTotal += 1
      }
    }
  }
  const coveragePct = trackedTotal === 0 ? 0 : (presentTotal / trackedTotal) * 100

  // powers the toggle so flipping to dimensions mode also flips the headline number.
  let presentDim = 0
  let trackedDim = 0
  for (const r of rows) {
    const result = reclassifyRow(r, COLUMNS, classifyByCategory)
    presentDim += result.coverage
    trackedDim += result.expected
  }
  const coveragePctDimensions = trackedDim === 0 ? 0 : (presentDim / trackedDim) * 100

  const updatedAt =
    typeof raw.generatedAt === 'string'
      ? raw.generatedAt
      : statSync(RESOLVED_SNAPSHOT_PATH).mtime.toISOString()
  const stats: SnapshotStats = {
    tracked: rows.length,
    coveragePct,
    coveragePctDimensions,
    updatedAt,
  }

  return { columns, rows, total: rows.length, stats }
}

function readRaw(): RawSnapshot {
  let parsed: unknown
  try {
    parsed = JSON.parse(readFileSync(RESOLVED_SNAPSHOT_PATH, 'utf8'))
  } catch (err) {
    if ((err as NodeJS.ErrnoException).code === 'ENOENT') {
      throw new Error(
        `snapshot not found at ${RESOLVED_SNAPSHOT_PATH}; run scripts/refresh.sh to produce it`,
      )
    }
    throw err
  }
  if (
    parsed == null ||
    typeof parsed !== 'object' ||
    !Array.isArray((parsed as RawSnapshot).cells) ||
    !Array.isArray((parsed as RawSnapshot).protocols)
  ) {
    throw new Error(
      `snapshot at ${RESOLVED_SNAPSHOT_PATH} is malformed; expected {cells:[], protocols:[]}`,
    )
  }
  return parsed as RawSnapshot
}

function presenceBySlug(cells: RawCell[]): Map<string, Set<string>> {
  const out = new Map<string, Set<string>>()
  for (const c of cells) {
    let set = out.get(c.slug)
    if (!set) {
      set = new Set()
      out.set(c.slug, set)
    }
    set.add(c.metric)
  }
  return out
}

export function projectRow(p: RawProtocol, present: Set<string> | undefined): Row {
  const cells = {} as Cells
  let coverage = 0
  let expected = 0
  for (const col of COLUMNS) {
    const isPresent = present !== undefined && present.has(col.key)
    const state = classifyCell(p.dimTypes, col.key, isPresent)
    cells[col.key] = state
    if (state === 'present') coverage += 1
    if (state === 'present' || state === 'missing') expected += 1
  }
  // Empty-string category from the wire collapses to undefined so the
  // category filter and chip strip do not show a blank entry.
  const category = p.category != null && p.category !== '' ? p.category : undefined
  // Lowercase defensively so the chain filter matches its URL token regardless
  // of upstream casing drift; today the upstream already emits lowercase.
  const chains = p.chains.map((c) => c.toLowerCase())
  return {
    slug: p.slug,
    name: p.name,
    category,
    chains,
    cells,
    coverage,
    expected,
    dimTypes: p.dimTypes,
  }
}
