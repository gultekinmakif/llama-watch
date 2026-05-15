// Package dimensions parses the upstream DefiLlama clones into the db tables
package dimensions

import "strings"

// Canonical mirrors defillama-server/defi/src/utils/sluggify.ts:3 (sluggifyString)
func Canonical(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "'", "")
	return s
}

// FromFilename is the default fallback.
// It is used when `data{N}.ts` did not ref the adapter to a TVL adapter.
func FromFilename(basename string) string {
	s := strings.ToLower(basename)
	s = strings.TrimSuffix(s, ".ts")
	s = strings.TrimSuffix(s, ".js")
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
		s = strings.ReplaceAll(s, "'", "")
	}
	s = strings.Trim(s, "-")
	return s
}
