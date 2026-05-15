package dimensions

import (
	"reflect"
	"sort"
	"testing"

	"github.com/gultekinmakif/llama-watch/internal/testutil"
)

func relset(cs []Adapter) []string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.RelPath)
	}
	sort.Strings(out)
	return out
}

func indexByRel(cs []Adapter) map[string]Adapter {
	out := make(map[string]Adapter, len(cs))
	for _, c := range cs {
		out[c.RelPath] = c
	}
	return out
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
	testutil.WriteFile(t, root, "dimension-adapters/fees/wbtc.ts", "")
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
	testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/uniswap-v2/index.js", "")
	// projects/<slug>.ts
	testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/wbtc.ts", "")
	// projects/<slug>.js (bare flat .js shape with no parent dir)
	testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/foo.js", "")
	// projects/<slug>/index.ts
	testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/aave-v2/index.ts", "")

	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 4 {
		t.Fatalf("want 4 candidates, got %d (%v)", len(got), relset(got))
	}
	byRel := indexByRel(got)

	cases := []struct {
		rel  string
		slug string
	}{
		{"DefiLlama-Adapters/projects/uniswap-v2/index.js", "uniswap-v2"},
		{"DefiLlama-Adapters/projects/wbtc.ts", "wbtc"},
		{"DefiLlama-Adapters/projects/foo.js", "foo"},
		{"DefiLlama-Adapters/projects/aave-v2/index.ts", "aave-v2"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.rel, func(t *testing.T) {
			c, ok := byRel[tc.rel]
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
	testutil.WriteFile(t, root, "dimension-adapters/fees/aave-v2.ts", "")
	testutil.WriteFile(t, root, "dimension-adapters/fees/wbtc.ts", "")
	testutil.WriteFile(t, root, "dimension-adapters/dexs/uniswap-v3/index.ts", "")
	testutil.WriteFile(t, root, "dimension-adapters/options/lyra.ts", "")
	testutil.WriteFile(t, root, "dimension-adapters/open-interest/gmx/index.ts", "")

	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 5 {
		t.Fatalf("want 5 candidates, got %d (%v)", len(got), relset(got))
	}
	byRel := indexByRel(got)

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
			c, ok := byRel[tc.rel]
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
	// Multi-version protocol siblings are separate rows and must not collapse. Uniswap V2, V3, V4 each get their own slug.
	root := t.TempDir()
	testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/uniswap-v2/index.js", "")
	testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/uniswap-v3/index.js", "")
	testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/uniswap-v4/index.js", "")

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
	t.Run("hidden_dirs", func(t *testing.T) {
		root := t.TempDir()
		testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/.git/HEAD", "")
		testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/.github/workflows/ci.yml", "")
		testutil.WriteFile(t, root, "dimension-adapters/.git/HEAD", "")
		testutil.WriteFile(t, root, "dimension-adapters/fees/.hidden/sneaky.ts", "")
		got, err := Walk(root)
		if err != nil {
			t.Fatalf("Walk: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("want 0, got %v", relset(got))
		}
	})

	t.Run("node_modules", func(t *testing.T) {
		root := t.TempDir()
		testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/foo/node_modules/lib/a.ts", "")
		testutil.WriteFile(t, root, "dimension-adapters/fees/node_modules/dep/index.ts", "")
		got, err := Walk(root)
		if err != nil {
			t.Fatalf("Walk: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("want 0, got %v", relset(got))
		}
	})

	t.Run("declaration_files_dot_d_ts", func(t *testing.T) {
		root := t.TempDir()
		testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/types.d.ts", "")
		testutil.WriteFile(t, root, "dimension-adapters/fees/types.d.ts", "")
		got, err := Walk(root)
		if err != nil {
			t.Fatalf("Walk: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("want 0, got %v", relset(got))
		}
	})

	t.Run("excluded_tvl_subtrees", func(t *testing.T) {
		root := t.TempDir()
		for _, excluded := range []string{"helper", "treasury", "entities", "config", "stacks", "test"} {
			testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/"+excluded+"/util.ts", "")
		}
		got, err := Walk(root)
		if err != nil {
			t.Fatalf("Walk: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("want 0, got %v", relset(got))
		}
	})

	t.Run("non_js_ts_files", func(t *testing.T) {
		root := t.TempDir()
		testutil.WriteFile(t, root, "dimension-adapters/fees/README.md", "")
		testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/foo/package.json", "")
		got, err := Walk(root)
		if err != nil {
			t.Fatalf("Walk: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("want 0, got %v", relset(got))
		}
	})

	t.Run("real_candidates_survive_mixed_noise", func(t *testing.T) {
		root := t.TempDir()
		// Real adapters that must be picked up.
		testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/wbtc.ts", "")
		testutil.WriteFile(t, root, "dimension-adapters/fees/aave-v2.ts", "")
		// One representative of every skip category, intermixed.
		testutil.WriteFile(t, root, "DefiLlama-Adapters/.git/HEAD", "")
		testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/.github/workflows/ci.yml", "")
		testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/foo/node_modules/lib/a.ts", "")
		testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/types.d.ts", "")
		testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/helper/util.ts", "")
		testutil.WriteFile(t, root, "dimension-adapters/fees/README.md", "")

		got, err := Walk(root)
		if err != nil {
			t.Fatalf("Walk: %v", err)
		}

		wantRels := []string{
			"DefiLlama-Adapters/projects/wbtc.ts",
			"dimension-adapters/fees/aave-v2.ts",
		}
		if g := relset(got); !reflect.DeepEqual(g, wantRels) {
			t.Fatalf("want %v, got %v", wantRels, g)
		}

		byRel := indexByRel(got)
		if c := byRel["DefiLlama-Adapters/projects/wbtc.ts"]; c.Type != "tvl" || c.Slug != "wbtc" {
			t.Fatalf("tvl candidate wrong: %+v", c)
		}
		if c := byRel["dimension-adapters/fees/aave-v2.ts"]; c.Type != "fees" || c.Slug != "aave-v2" {
			t.Fatalf("dim candidate wrong: %+v", c)
		}
	})
}

func TestWalk_FilenameNormalization(t *testing.T) {
	// FromFilename is exercised here through resolveSlug, mainly to confirm
	// underscores and casing in directory/file names route through correctly.
	root := t.TempDir()
	testutil.WriteFile(t, root, "DefiLlama-Adapters/projects/Uniswap_V2/index.js", "")
	testutil.WriteFile(t, root, "dimension-adapters/fees/Uniswap_V2.ts", "")
	testutil.WriteFile(t, root, "dimension-adapters/dexs/Foo_Bar/index.ts", "")

	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("want 3, got %d (%v)", len(got), relset(got))
	}
	byRel := indexByRel(got)

	want := map[string]string{
		"DefiLlama-Adapters/projects/Uniswap_V2/index.js": "uniswap-v2",
		"dimension-adapters/fees/Uniswap_V2.ts":           "uniswap-v2",
		"dimension-adapters/dexs/Foo_Bar/index.ts":        "foo-bar",
	}
	for rel, slug := range want {
		c, ok := byRel[rel]
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
	testutil.Mkdir(t, root, "DefiLlama-Adapters/projects")
	testutil.Mkdir(t, root, "dimension-adapters/fees")
	testutil.Mkdir(t, root, "dimension-adapters/dexs")

	got, err := Walk(root)
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("want 0, got %d (%v)", len(got), relset(got))
	}
}
