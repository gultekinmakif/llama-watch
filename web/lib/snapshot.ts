// Build-time snapshot reader. Projects the wire shape {cells, protocols}
// onto the closed column set the table renders. Keep COLUMNS in lockstep
// with internal/registry/columns.go.

import { readFileSync } from 'node:fs'
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

export type CellState = 'na' | 'missing' | 'present' | 'over'

// Mirror of internal/registry/expectations.go. Keep in lockstep with the Go seed.
// Drift between this map and Go means the build-time matrix disagrees with /api/matrix.
const EXPECTATIONS: Record<string, Record<string, true>> = {
  Lending: { tvl: true, dailyFees: true, dailyRevenue: true, dailyHoldersRevenue: true, dailyUserFees: true, dailySupplySideRevenue: true, dailyProtocolRevenue: true },
  'DEX Aggregator': { dailyVolume: true, dailyFees: true },
  Derivatives: { tvl: true, dailyVolume: true, openInterestAtEnd: true, longOpenInterestAtEnd: true, shortOpenInterestAtEnd: true, dailyFees: true, dailyRevenue: true, dailyHoldersRevenue: true },
  Options: { tvl: true, dailyNotionalVolume: true, dailyPremiumVolume: true, dailyFees: true },
  Bridge: { tvl: true, dailyBridgeVolume: true, dailyFees: true },
  'Canonical Bridge': { tvl: true, dailyBridgeVolume: true },
  'Cross Chain Bridge': { tvl: true, dailyBridgeVolume: true, dailyFees: true },
  'Bridge Aggregator': { dailyBridgeVolume: true, dailyVolume: true, dailyFees: true },
  Insurance: { tvl: true, dailyFees: true, dailyRevenue: true, dailyHoldersRevenue: true },
  'Liquid Staking': { tvl: true, dailyFees: true, dailyRevenue: true, dailyHoldersRevenue: true, dailySupplySideRevenue: true },
  CDP: { tvl: true, dailyFees: true, dailyRevenue: true, dailyHoldersRevenue: true, dailyUserFees: true },
  'CDP Manager': { tvl: true, dailyFees: true, dailyRevenue: true },
  Synthetics: { tvl: true, dailyVolume: true, dailyFees: true, dailyRevenue: true },
  Yield: { tvl: true, dailyFees: true, dailyRevenue: true },
  'Yield Aggregator': { tvl: true, dailyFees: true, dailyRevenue: true },
  Farm: { tvl: true, dailyFees: true, dailyRevenue: true },
  'Leveraged Farming': { tvl: true, dailyFees: true, dailyRevenue: true },
  'Algo-Stables': { tvl: true, dailyFees: true },
  Chain: { tvl: true, dailyFees: true, dailyRevenue: true, dailyTransactionsCount: true, dailyGasUsed: true, dailyActiveUsers: true, dailyNewUsers: true },
  Indexes: { tvl: true, dailyFees: true, dailyRevenue: true },
}

// Classifies a (category, metric, present) triple into one of na/missing/present/over.
// Mirrors internal/registry/expectations.go::ClassifyCell. Unseeded categories fall
// through: present is CellPresent, absent is CellNA.
export function classifyCell(category: string, metric: string, present: boolean): CellState {
  const seed = EXPECTATIONS[category]
  if (!seed) return present ? 'present' : 'na'
  const expected = seed[metric] === true
  if (present && expected) return 'present'
  if (present && !expected) return 'over'
  if (!present && expected) return 'missing'
  return 'na'
}

export type Cells = Record<ColumnKey, CellState>

export interface Row {
  slug: string
  name: string
  category?: string
  chains: string[]
  cells: Cells
  // Precomputed at build time so sort/render does not redo Object.values on every pass.
  coverage: number
}

export interface Snapshot {
  columns: Column[]
  rows: Row[]
  total: number
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
}

export interface RawSnapshot {
  cells: RawCell[]
  protocols: RawProtocol[]
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
  return { columns: COLUMNS.map((c) => ({ ...c })), rows, total: rows.length }
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
  for (const col of COLUMNS) {
    const isPresent = present !== undefined && present.has(col.key)
    const state = classifyCell(p.category ?? '', col.key, isPresent)
    cells[col.key] = state
    if (state === 'present') coverage += 1
  }
  // Empty-string category from the wire collapses to undefined so the
  // category filter and chip strip do not show a blank entry.
  const category = p.category != null && p.category !== '' ? p.category : undefined
  // Lowercase defensively so the chain filter matches its URL token regardless
  // of upstream casing drift; today the upstream already emits lowercase.
  const chains = p.chains.map((c) => c.toLowerCase())
  return { slug: p.slug, name: p.name, category, chains, cells, coverage }
}
