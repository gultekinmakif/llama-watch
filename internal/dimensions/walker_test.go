package dimensions

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// mkfile writes an empty file at rel under root, creating parents.
// The walker never reads file contents, so content-less fixtures are fine.
func mkfile(t *testing.T, root, rel string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(full), err)
	}
	f, err := os.Create(full)
	if err != nil {
		t.Fatalf("create %s: %v", full, err)
	}
	f.Close()
}

func mkdir(t *testing.T, root, rel string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(full, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", full, err)
	}
}

func relset(cs []Adapter) []string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.RelPath)
	}
	sort.Strings(out)
	return out
}

func findByRel(cs []Adapter, rel string) (Adapter, bool) {
	for _, c := range cs {
		if c.RelPath == rel {
			return c, true
		}
	}
	return Adapter{}, false
}

func TestWalk_MissingUpstreamDirs(t *testing.T) {
	root := t.TempDir()
	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk on empty dir: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("want 0 candidates, got %d (%v)", len(got), relset(got))
	}
}

func TestWalk_PartialMissing(t *testing.T) {
	root := t.TempDir()
	// Only one of the two repos exists.
	mkfile(t, root, "dimension-adapters/fees/wbtc.ts")
	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 candidate, got %d (%v)", len(got), relset(got))
	}
	if got[0].Type != "fees" {
		t.Fatalf("unexpected candidate: %+v", got[0])
	}
}

func TestWalk_TVLShapes(t *testing.T) {
	root := t.TempDir()
	// projects/<slug>/index.js
	mkfile(t, root, "DefiLlama-Adapters/projects/uniswap-v2/index.js")
	// projects/<slug>.ts
	mkfile(t, root, "DefiLlama-Adapters/projects/wbtc.ts")
	// projects/<slug>/index.ts
	mkfile(t, root, "DefiLlama-Adapters/projects/aave-v2/index.ts")

	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3 candidates, got %d (%v)", len(got), relset(got))
	}

	cases := []struct {
		rel  string
		slug string
	}{
		{"DefiLlama-Adapters/projects/uniswap-v2/index.js", "uniswap-v2"},
		{"DefiLlama-Adapters/projects/wbtc.ts", "wbtc"},
		{"DefiLlama-Adapters/projects/aave-v2/index.ts", "aave-v2"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.rel, func(t *testing.T) {
			c, ok := findByRel(got, tc.rel)
			if !ok {
				t.Fatalf("missing candidate %s in %v", tc.rel, relset(got))
			}
			if c.Type != "tvl" {
				t.Fatalf("type = %q, want tvl", c.Type)
			}
			if c.Slug != tc.slug {
				t.Fatalf("slug = %q, want %q", c.Slug, tc.slug)
			}
			if c.AbsPath == "" {
				t.Fatalf("abs path empty")
			}
		})
	}
}

func TestWalk_DimensionShapes(t *testing.T) {
	root := t.TempDir()
	mkfile(t, root, "dimension-adapters/fees/aave-v2.ts")
	mkfile(t, root, "dimension-adapters/fees/wbtc.ts")
	mkfile(t, root, "dimension-adapters/dexs/uniswap-v3/index.ts")
	mkfile(t, root, "dimension-adapters/options/lyra.ts")
	mkfile(t, root, "dimension-adapters/open-interest/gmx/index.ts")

	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 5 {
		t.Fatalf("want 5 candidates, got %d (%v)", len(got), relset(got))
	}

	cases := []struct {
		rel  string
		typ  string
		slug string
	}{
		{"dimension-adapters/fees/aave-v2.ts", "fees", "aave-v2"},
		{"dimension-adapters/fees/wbtc.ts", "fees", "wbtc"},
		{"dimension-adapters/dexs/uniswap-v3/index.ts", "dexs", "uniswap-v3"},
		{"dimension-adapters/options/lyra.ts", "options", "lyra"},
		{"dimension-adapters/open-interest/gmx/index.ts", "open-interest", "gmx"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.rel, func(t *testing.T) {
			c, ok := findByRel(got, tc.rel)
			if !ok {
				t.Fatalf("missing candidate %s in %v", tc.rel, relset(got))
			}
			if c.Type != tc.typ {
				t.Fatalf("type = %q, want %q", c.Type, tc.typ)
			}
			if c.Slug != tc.slug {
				t.Fatalf("slug = %q, want %q", c.Slug, tc.slug)
			}
		})
	}
}

func TestWalk_MultiVersionSiblings(t *testing.T) {
	// PARSER.md §132 — must NOT collapse Uniswap V2/V3/V4.
	root := t.TempDir()
	mkfile(t, root, "DefiLlama-Adapters/projects/uniswap-v2/index.js")
	mkfile(t, root, "DefiLlama-Adapters/projects/uniswap-v3/index.js")
	mkfile(t, root, "DefiLlama-Adapters/projects/uniswap-v4/index.js")

	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3, got %d", len(got))
	}
	seen := map[string]bool{}
	for _, c := range got {
		seen[c.Slug] = true
	}
	for _, want := range []string{"uniswap-v2", "uniswap-v3", "uniswap-v4"} {
		if !seen[want] {
			t.Fatalf("missing %s; got %v", want, seen)
		}
	}
}

func TestWalk_Skips(t *testing.T) {
	root := t.TempDir()

	// Real adapter that must be picked up.
	mkfile(t, root, "DefiLlama-Adapters/projects/wbtc.ts")
	mkfile(t, root, "dimension-adapters/fees/aave-v2.ts")

	// Hidden dirs anywhere in the tree.
	mkfile(t, root, "DefiLlama-Adapters/.git/HEAD")
	mkfile(t, root, "DefiLlama-Adapters/projects/.github/workflows/ci.yml")
	mkfile(t, root, "dimension-adapters/.git/HEAD")
	mkfile(t, root, "dimension-adapters/fees/.hidden/sneaky.ts")

	// node_modules anywhere.
	mkfile(t, root, "DefiLlama-Adapters/projects/foo/node_modules/lib/a.ts")
	mkfile(t, root, "dimension-adapters/fees/node_modules/dep/index.ts")

	// .d.ts declaration files.
	mkfile(t, root, "DefiLlama-Adapters/projects/types.d.ts")
	mkfile(t, root, "dimension-adapters/fees/types.d.ts")

	// PARSER.md §21 excluded subtrees under projects/.
	for _, excluded := range []string{"helper", "treasury", "entities", "config", "stacks", "test"} {
		mkfile(t, root, "DefiLlama-Adapters/projects/"+excluded+"/util.ts")
	}

	// Non-JS/TS noise.
	mkfile(t, root, "dimension-adapters/fees/README.md")
	mkfile(t, root, "DefiLlama-Adapters/projects/foo/package.json")

	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}

	wantRels := []string{
		"DefiLlama-Adapters/projects/wbtc.ts",
		"dimension-adapters/fees/aave-v2.ts",
	}
	if g := relset(got); !equalStrs(g, wantRels) {
		t.Fatalf("want %v, got %v", wantRels, g)
	}
}

func TestWalk_FilenameNormalization(t *testing.T) {
	// FromFilename is exercised here through resolveSlug, mainly to confirm
	// underscores and casing in directory/file names route through correctly.
	root := t.TempDir()
	mkfile(t, root, "DefiLlama-Adapters/projects/Uniswap_V2/index.js")
	mkfile(t, root, "dimension-adapters/fees/Uniswap_V2.ts")
	mkfile(t, root, "dimension-adapters/dexs/Foo_Bar/index.ts")

	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3, got %d (%v)", len(got), relset(got))
	}

	want := map[string]string{
		"DefiLlama-Adapters/projects/Uniswap_V2/index.js": "uniswap-v2",
		"dimension-adapters/fees/Uniswap_V2.ts":           "uniswap-v2",
		"dimension-adapters/dexs/Foo_Bar/index.ts":        "foo-bar",
	}
	for rel, slug := range want {
		c, ok := findByRel(got, rel)
		if !ok {
			t.Fatalf("missing %s", rel)
		}
		if c.Slug != slug {
			t.Fatalf("%s: slug = %q, want %q", rel, c.Slug, slug)
		}
	}
}

func TestWalk_EmptyDirsNoCandidates(t *testing.T) {
	root := t.TempDir()
	mkdir(t, root, "DefiLlama-Adapters/projects")
	mkdir(t, root, "dimension-adapters/fees")
	mkdir(t, root, "dimension-adapters/dexs")

	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("want 0, got %d (%v)", len(got), relset(got))
	}
}

func equalStrs(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
