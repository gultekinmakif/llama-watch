// Build-time snapshot reader.

import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

const RESOLVED_SNAPSHOT_PATH = resolve(
  process.cwd(),
  '..',
  'var',
  'snapshot',
  'snapshot.json',
)

export type ColumnKey =
  | 'tvl'
  | 'dailyFees'
  | 'dailyRevenue'
  | 'dailyVolume'
  | 'dailyNotionalVolume'
  | 'dailyPremiumVolume'
  | 'openInterestAtEnd'
  | 'dailyBridgeVolume'
  | 'dailyActiveUsers'

export interface Column {
  key: ColumnKey
  label: string
}

// 1 = an adapter file exists for this (protocol, dimension)
// 0 = no adapter file.
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

// Runs at build time only; sync fs is correct here.
export function loadSnapshot(): Snapshot {
  let raw: string
  try {
    raw = readFileSync(RESOLVED_SNAPSHOT_PATH, 'utf8')
  } catch (err) {
    if ((err as NodeJS.ErrnoException).code === 'ENOENT') {
      throw new Error(
        `snapshot not found at ${RESOLVED_SNAPSHOT_PATH}; run scripts/refresh.sh to produce it`,
      )
    }
    throw err
  }
  return JSON.parse(raw) as unknown as Snapshot
}
