// 1.5: Detect sub-metric keys in a dimension-adapter source via word-boundary regex.
// TODO: v1 tolerates false positives from comments and strings; the matrix is presence/absence.
package dimensions

import (
	"os"
	"regexp"

	"github.com/gultekinmakif/llama-watch/internal/registry"
)

// detectionRx is the word-boundary regex per kind, compiled.
var detectionRx = func() map[string]*regexp.Regexp {
	seen := map[string]*regexp.Regexp{}
	for _, kinds := range registry.AllMetricsByDimension() {
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
	cands := registry.MetricsFor(dimType)
	if cands == nil {
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
