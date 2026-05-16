// 1.6: Slug canonicalization helpers used throughout stage 1.
// dimensions package: joins upstream adapter clones with the extracted protocols JSON.
package dimensions

import (
	"path"
	"strings"
)

// Canonical mirrors defillama-server/defi/src/utils/sluggify.ts:3 (sluggifyString)
func Canonical(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "'", "")
	return s
}

// pathSlug extracts the canonical slug from a slash-form path. If the
func pathSlug(p string) string {
	p = strings.TrimRight(p, "/")
	if p == "" {
		return ""
	}
	base := path.Base(p)
	if strings.TrimSuffix(strings.TrimSuffix(base, ".ts"), ".js") == "index" {
		if parent := path.Dir(p); parent != "." && parent != "/" {
			return FromFilename(path.Base(parent))
		}
	}
	return FromFilename(base)
}

// FromFilename is the default fallback.
// It is used when `data{N}.ts` did not ref the adapter to a TVL adapter.
//
// > Intentionally diverges from Canonical: no `'` stripping, because
// filenames don't carry apostrophes and there are protocols with `'` in their names.
func FromFilename(basename string) string {
	s := strings.ToLower(basename)
	s = strings.TrimSuffix(s, ".ts")
	s = strings.TrimSuffix(s, ".js")
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	return s
}
