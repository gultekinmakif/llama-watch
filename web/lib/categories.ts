// Category taxonomy, web-only. Go and build-snapshot only need presets.json.

import type { ColumnKey } from './snapshot'

export const CATEGORIES_EXPECTED: Record<string, readonly ColumnKey[]> = {
  "Algo-Stables": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyNotionalVolume", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Bridge": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyNotionalVolume", "dailyBridgeVolume", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Bridge Aggregator": ["tvl", "dailyFees", "dailyVolume", "dailyBridgeVolume"],
  "CDP": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyNotionalVolume", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "CDP Manager": ["tvl", "dailyFees", "dailyRevenue", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Canonical Bridge": ["tvl", "dailyBridgeVolume"],
  "Chain": ["tvl", "dailyFees", "dailyRevenue", "dailyActiveUsers", "dailyTransactionsCount", "dailyGasUsed", "dailyNewUsers"],
  "Cross Chain Bridge": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyNotionalVolume", "dailyBridgeVolume", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "DEX Aggregator": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyNotionalVolume", "dailyBridgeVolume", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Derivatives": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyNotionalVolume", "dailyPremiumVolume", "openInterestAtEnd", "dailyActiveUsers", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes", "longOpenInterestAtEnd", "shortOpenInterestAtEnd", "dailyTransactionsCount", "dailyGasUsed", "dailyNormalizedVolume", "dailyActiveLiquidity"],
  "Farm": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyNotionalVolume", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Indexes": ["tvl", "dailyFees", "dailyRevenue", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Insurance": ["tvl", "dailyFees", "dailyRevenue", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Lending": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyNotionalVolume", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Leveraged Farming": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyNotionalVolume", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Liquid Staking": ["tvl", "dailyFees", "dailyRevenue", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Options": ["tvl", "dailyFees", "dailyRevenue", "dailyNotionalVolume", "dailyPremiumVolume", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Synthetics": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyNotionalVolume", "dailyPremiumVolume", "openInterestAtEnd", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes", "longOpenInterestAtEnd", "shortOpenInterestAtEnd"],
  "Yield": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyNotionalVolume", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
  "Yield Aggregator": ["tvl", "dailyFees", "dailyRevenue", "dailyVolume", "dailyUserFees", "dailyHoldersRevenue", "dailyProtocolRevenue", "dailySupplySideRevenue", "dailyCreatorRevenue", "dailyBribesRevenue", "dailyTokenTaxes"],
}
