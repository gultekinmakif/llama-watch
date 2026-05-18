// Pure-runtime helpers backing the category and dimType filter dropdowns.
// No node imports so client components can consume these directly.

import { EXPECTATIONS } from './cell-state'
import type { ColumnKey } from './snapshot'

// dimType to metric column keys, ported from tools/build-snapshot.ts:KEYS_TO_STORE.
// Keep in lockstep with the Go and TS sides; the upstream is defillama-server getDimensionsConfig.
const KEYS_TO_STORE = {
  dexs: ['dailyVolume', 'dailyNotionalVolume'],
  derivatives: ['dailyVolume'],
  aggregators: ['dailyVolume'],
  'aggregator-derivatives': ['dailyVolume'],
  fees: [
    'dailyFees',
    'dailyRevenue',
    'dailyUserFees',
    'dailySupplySideRevenue',
    'dailyProtocolRevenue',
    'dailyHoldersRevenue',
    'dailyCreatorRevenue',
    'dailyBribesRevenue',
    'dailyTokenTaxes',
  ],
  options: ['dailyPremiumVolume', 'dailyNotionalVolume'],
  'open-interest': ['openInterestAtEnd', 'shortOpenInterestAtEnd', 'longOpenInterestAtEnd'],
  'bridge-aggregators': ['dailyBridgeVolume'],
  'active-users': ['dailyActiveUsers', 'dailyTransactionsCount', 'dailyGasUsed'],
  'new-users': ['dailyNewUsers'],
  'nft-volume': ['dailyVolume'],
  'normalized-volume': ['dailyNormalizedVolume', 'dailyActiveLiquidity'],
  incentives: ['tokenIncentives'],
} satisfies Record<string, readonly ColumnKey[]>

export const CATEGORIES: readonly string[] = Object.keys(EXPECTATIONS).sort()

export const DIMTYPES: readonly string[] = Object.keys(KEYS_TO_STORE).sort()

// Returns a fresh array so callers can mutate without poisoning the seed.
export function expectedColumnsFor(category: string): ColumnKey[] {
  const seed = EXPECTATIONS[category]
  if (!seed) return []
  return Object.keys(seed) as ColumnKey[]
}

// Returns a fresh array so callers can mutate without poisoning the port.
export function metricsFor(dimType: string): ColumnKey[] {
  const metrics = (KEYS_TO_STORE as Record<string, readonly ColumnKey[]>)[dimType]
  if (!metrics) return []
  return [...metrics]
}
