// Pinned matrix column set shared by the /api/matrix handler and snapshot writer.
package registry

type Column struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

// Closed, fixed column set; order is load-bearing. Use Columns() for a safe copy.
// Keys mirror defillama-server getDimensionsConfig KEYS_TO_STORE flattened across dimTypes, plus tvl.
var columns = []Column{
	{Key: "tvl", Label: "TVL"},
	{Key: "dailyFees", Label: "Daily Fees"},
	{Key: "dailyRevenue", Label: "Daily Revenue"},
	{Key: "dailyVolume", Label: "Daily Volume"},
	{Key: "dailyNotionalVolume", Label: "Notional Volume"},
	{Key: "dailyPremiumVolume", Label: "Premium Volume"},
	{Key: "openInterestAtEnd", Label: "Open Interest"},
	{Key: "dailyBridgeVolume", Label: "Bridge Volume"},
	{Key: "dailyActiveUsers", Label: "Active Users"},
	{Key: "dailyUserFees", Label: "User Fees"},
	{Key: "dailyHoldersRevenue", Label: "Holders Revenue"},
	{Key: "dailyProtocolRevenue", Label: "Protocol Revenue"},
	{Key: "dailySupplySideRevenue", Label: "Supply-Side Revenue"},
	{Key: "dailyCreatorRevenue", Label: "Creator Revenue"},
	{Key: "dailyBribesRevenue", Label: "Bribes Revenue"},
	{Key: "dailyTokenTaxes", Label: "Token Taxes"},
	{Key: "longOpenInterestAtEnd", Label: "Long Open Interest"},
	{Key: "shortOpenInterestAtEnd", Label: "Short Open Interest"},
	{Key: "dailyTransactionsCount", Label: "Transactions"},
	{Key: "dailyGasUsed", Label: "Gas Used"},
	{Key: "dailyNewUsers", Label: "New Users"},
	{Key: "dailyNormalizedVolume", Label: "Normalized Volume"},
	{Key: "dailyActiveLiquidity", Label: "Active Liquidity"},
	{Key: "tokenIncentives", Label: "Token Incentives"},
}

// Columns returns a copy of the pinned column set so external callers can read
// the matrix shape without risking mutation of the package-level slice.
func Columns() []Column {
	out := make([]Column, len(columns))
	copy(out, columns)
	return out
}
