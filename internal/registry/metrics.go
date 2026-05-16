// Dimension-to-kinds seed for the matrix.
package registry

// metricsByDimension maps a dimension-adapter directory to candidate sub-metric kinds per file.
var metricsByDimension = map[string][]string{
	"fees":                   {"dailyFees", "dailyRevenue"},
	"options":                {"dailyNotionalVolume", "dailyPremiumVolume"},
	"aggregator-options":     {"dailyNotionalVolume", "dailyPremiumVolume"},
	"dexs":                   {"dailyVolume"},
	"aggregators":            {"dailyVolume"},
	"aggregator-derivatives": {"dailyVolume"},
	"derivatives":            {"dailyVolume"},
	"bridge-aggregators":     {"dailyBridgeVolume"},
	"open-interest":          {"openInterestAtEnd"},
	"active-users":           {"dailyActiveUsers"},
	"users":                  {"dailyActiveUsers"},
}

// MetricsFor returns a copy of the kinds slice for dimType; nil if unknown.
func MetricsFor(dimType string) []string {
	kinds, ok := metricsByDimension[dimType]
	if !ok {
		return nil
	}
	out := make([]string, len(kinds))
	copy(out, kinds)
	return out
}

// AllMetricsByDimension returns a fresh copy of the full map with every inner slice copied.
func AllMetricsByDimension() map[string][]string {
	out := make(map[string][]string, len(metricsByDimension))
	for dim, kinds := range metricsByDimension {
		dup := make([]string, len(kinds))
		copy(dup, kinds)
		out[dim] = dup
	}
	return out
}
