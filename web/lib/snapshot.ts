// Build-time snapshot reader. Projects the wire shape {cells, protocols}
// onto the closed column set the table renders. Keep COLUMNS in lockstep
// with internal/registry/columns.go.

import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

const RESOLVED_SNAPSHOT_PATH = resolve(
  process.cwd(),
  '..',
  'var',
  'snapshot',
  'snapshot.json',
)

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

// 1 = an adapter file exists for this (protocol, metric); 0 = absent.
export type Cells = Record<ColumnKey, 0 | 1>

export interface Row {
  slug: string
  name: string
  category?: string
  chains: string[]
  cells: Cells
}

export interface Snapshot {
  columns: Column[]
  rows: Row[]
  total: number
}

interface RawCell {
  slug: string
  metric: string
  codePath: string
}

interface RawProtocol {
  slug: string
  name: string
  category: string
  chains: string[]
  dataFile: string
}

interface RawSnapshot {
  cells: RawCell[]
  protocols: RawProtocol[]
}

// Runs at build time only; sync fs is correct here.
export function loadSnapshot(): Snapshot {
  const raw = readRaw()
  const presence = presenceBySlug(raw.cells)
  const rows = raw.protocols.map((p) => projectRow(p, presence.get(p.slug)))
  return { columns: COLUMNS.map((c) => ({ ...c })), rows, total: rows.length }
}

function readRaw(): RawSnapshot {
  try {
    return JSON.parse(readFileSync(RESOLVED_SNAPSHOT_PATH, 'utf8')) as RawSnapshot
  } catch (err) {
    if ((err as NodeJS.ErrnoException).code === 'ENOENT') {
      throw new Error(
        `snapshot not found at ${RESOLVED_SNAPSHOT_PATH}; run scripts/refresh.sh to produce it`,
      )
    }
    throw err
  }
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

function projectRow(p: RawProtocol, present: Set<string> | undefined): Row {
  const cells = {} as Cells
  for (const col of COLUMNS) {
    cells[col.key] = present && present.has(col.key) ? 1 : 0
  }
  return { slug: p.slug, name: p.name, category: p.category, chains: p.chains, cells }
}
