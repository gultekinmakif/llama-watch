// kinds.go: detect sub-metric keys in a dimension-adapter source via word-boundary regex.
// TODO: v1 tolerates false positives from comments and strings; the matrix is presence/absence.
package dimensions

import (
	"os"
	"regexp"
)

// metricsLookup maps a dimension-adapter directory to candidate sub-metric kinds per file
var metricsLookup = map[string][]string{
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

// detectionRx is the word-boundary regex per kind, compiled.
var detectionRx = func() map[string]*regexp.Regexp {
	seen := map[string]*regexp.Regexp{}
	for _, kinds := range metricsLookup {
		for _, k := range kinds {
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = regexp.MustCompile(`\b` + regexp.QuoteMeta(k) + `\b`)
		}
	}
	return seen
}()

// DetectKinds returns the sub-metric kinds absPath exports. Unknown dimType returns nil.
func DetectKinds(absPath, dimType string) ([]string, error) {
	cands, ok := metricsLookup[dimType]
	if !ok {
		return nil, nil
	}
	contents, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, kind := range cands {
		if detectionRx[kind].Match(contents) {
			out = append(out, kind)
		}
	}
	return out, nil
}
