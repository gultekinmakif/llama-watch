// Pinned matrix column set shared by the /api/matrix handler and snapshot writer.
package registry

type Column struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

// Closed, fixed column set. Order is load-bearing. Unexported so callers cannot
// mutate the backing array; use Columns for a safe copy.
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
}

// Columns returns a copy of the pinned column set so external callers can read
// the matrix shape without risking mutation of the package-level slice.
func Columns() []Column {
	out := make([]Column, len(columns))
	copy(out, columns)
	return out
}
