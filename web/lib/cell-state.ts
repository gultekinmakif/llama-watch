// Pure-runtime four-state classifier. Lives separately from lib/snapshot.ts so
// client components can import it without pulling in node:fs / node:path / node:url.

export type CellState = 'na' | 'missing' | 'present' | 'over'

// Mirror of internal/registry/expectations.go.
// Drift between this map and Go means the build-time matrix disagrees with /api/matrix.
// Keep in lockstep with the Go seed.
export const EXPECTATIONS: Record<string, Record<string, true>> = {
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

// Classifies a (category, metric, present) triple. Mirrors internal/registry/expectations.go::ClassifyCell.
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
