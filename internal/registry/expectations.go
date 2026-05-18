// Category-to-expected-metrics seed and the four-state cell classifier.
package registry

import "maps"

// CellState is one of na, missing, present, over.
type CellState string

const (
	CellNA      CellState = "na"
	CellMissing CellState = "missing"
	CellPresent CellState = "present"
	CellOver    CellState = "over"
)

// Conservative seed; unseeded categories fall through in ClassifyCell.
var expectations = map[string]map[string]bool{
	"Lending": {
		"tvl":                    true,
		"dailyFees":              true,
		"dailyRevenue":           true,
		"dailyHoldersRevenue":    true,
		"dailyUserFees":          true,
		"dailySupplySideRevenue": true,
		"dailyProtocolRevenue":   true,
	},
	"DEX Aggregator": {
		"dailyVolume": true,
		"dailyFees":   true,
	},
	"Derivatives": {
		"tvl":                    true,
		"dailyVolume":            true,
		"openInterestAtEnd":      true,
		"longOpenInterestAtEnd":  true,
		"shortOpenInterestAtEnd": true,
		"dailyFees":              true,
		"dailyRevenue":           true,
		"dailyHoldersRevenue":    true,
	},
	"Options": {
		"tvl":                 true,
		"dailyNotionalVolume": true,
		"dailyPremiumVolume":  true,
		"dailyFees":           true,
	},
	"Bridge": {
		"tvl":               true,
		"dailyBridgeVolume": true,
		"dailyFees":         true,
	},
	"Canonical Bridge": {
		"tvl":               true,
		"dailyBridgeVolume": true,
	},
	"Cross Chain Bridge": {
		"tvl":               true,
		"dailyBridgeVolume": true,
		"dailyFees":         true,
	},
	"Bridge Aggregator": {
		"dailyBridgeVolume": true,
		"dailyVolume":       true,
		"dailyFees":         true,
	},
	"Insurance": {
		"tvl":                 true,
		"dailyFees":           true,
		"dailyRevenue":        true,
		"dailyHoldersRevenue": true,
	},
	"Liquid Staking": {
		"tvl":                    true,
		"dailyFees":              true,
		"dailyRevenue":           true,
		"dailyHoldersRevenue":    true,
		"dailySupplySideRevenue": true,
	},
	"CDP": {
		"tvl":                 true,
		"dailyFees":           true,
		"dailyRevenue":        true,
		"dailyHoldersRevenue": true,
		"dailyUserFees":       true,
	},
	"CDP Manager": {
		"tvl":          true,
		"dailyFees":    true,
		"dailyRevenue": true,
	},
	"Synthetics": {
		"tvl":          true,
		"dailyVolume":  true,
		"dailyFees":    true,
		"dailyRevenue": true,
	},
	"Yield": {
		"tvl":          true,
		"dailyFees":    true,
		"dailyRevenue": true,
	},
	"Yield Aggregator": {
		"tvl":          true,
		"dailyFees":    true,
		"dailyRevenue": true,
	},
	"Farm": {
		"tvl":          true,
		"dailyFees":    true,
		"dailyRevenue": true,
	},
	"Leveraged Farming": {
		"tvl":          true,
		"dailyFees":    true,
		"dailyRevenue": true,
	},
	"Algo-Stables": {
		"tvl":       true,
		"dailyFees": true,
	},
	"Chain": {
		"tvl":                    true,
		"dailyFees":              true,
		"dailyRevenue":           true,
		"dailyTransactionsCount": true,
		"dailyGasUsed":           true,
		"dailyActiveUsers":       true,
		"dailyNewUsers":          true,
	},
	"Indexes": {
		"tvl":          true,
		"dailyFees":    true,
		"dailyRevenue": true,
	},
}

// ExpectedMetrics returns the metrics this category should emit, or nil if unseeded.
// Returned map is a fresh copy.
func ExpectedMetrics(category string) map[string]bool {
	seed, ok := expectations[category]
	if !ok {
		return nil
	}
	return maps.Clone(seed)
}

// ClassifyCell returns the four-state coloring. Truth table:
//   present, expected     -> CellPresent      absent, expected     -> CellMissing
//   present, not expected -> CellOver         absent, not expected -> CellNA
// Unseeded categories: present -> CellPresent, absent -> CellNA.
func ClassifyCell(category, metric string, present bool) CellState {
	seed, ok := expectations[category]
	if !ok {
		if present {
			return CellPresent
		}
		return CellNA
	}
	expected := seed[metric]
	switch {
	case present && expected:
		return CellPresent
	case present && !expected:
		return CellOver
	case !present && expected:
		return CellMissing
	default:
		return CellNA
	}
}
