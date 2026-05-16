package dimensions

import (
	"reflect"
	"testing"
)

func TestModuleStem(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"uniswap-v2/index.js", "uniswap-v2"},
		{"wbtc.js", "wbtc"},
		{"wbtc.ts", "wbtc"},
		{"foo/index.ts", "foo"},
		{"bare-slug", "bare-slug"},
		{"", ""},
		{"a/b/index.js", "b"},
		{"trailing/", "trailing"},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			got := moduleStem(c.in)
			if got != c.want {
				t.Fatalf("moduleStem(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestDimensionCandidates(t *testing.T) {
	cases := []struct {
		name            string
		dimType, dimSlug string
		want            []string
	}{
		{
			name:    "common",
			dimType: "fees", dimSlug: "aave-v2",
			want: []string{
				"dimension-adapters/fees/aave-v2.ts",
				"dimension-adapters/fees/aave-v2/index.ts",
			},
		},
		{
			name:    "hyphenated dimType",
			dimType: "open-interest", dimSlug: "gmx",
			want: []string{
				"dimension-adapters/open-interest/gmx.ts",
				"dimension-adapters/open-interest/gmx/index.ts",
			},
		},
		{
			name:    "hyphenated slug",
			dimType: "dexs", dimSlug: "uniswap-v3",
			want: []string{
				"dimension-adapters/dexs/uniswap-v3.ts",
				"dimension-adapters/dexs/uniswap-v3/index.ts",
			},
		},
		{
			name:    "multi-version sibling slug",
			dimType: "fees", dimSlug: "uniswap-v2",
			want: []string{
				"dimension-adapters/fees/uniswap-v2.ts",
				"dimension-adapters/fees/uniswap-v2/index.ts",
			},
		},
		{
			name:    "empty inputs",
			dimType: "", dimSlug: "",
			want: []string{
				"dimension-adapters//.ts",
				"dimension-adapters///index.ts",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := dimensionCandidates(tc.dimType, tc.dimSlug)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("dimensionCandidates(%q, %q) = %#v, want %#v",
					tc.dimType, tc.dimSlug, got, tc.want)
			}
		})
	}
}

func TestIndexAdapters(t *testing.T) {
	in := []Adapter{
		{Type: "tvl", RelPath: "DefiLlama-Adapters/projects/a/index.js", AbsPath: "/x/a", Slug: "a"},
		{Type: "fees", RelPath: "dimension-adapters/fees/a.ts", AbsPath: "/x/fa", Slug: "a"},
	}
	idx := indexAdapters(in)
	if len(idx) != 2 {
		t.Fatalf("len = %d, want 2", len(idx))
	}
	if got := idx["dimension-adapters/fees/a.ts"]; got.Type != "fees" {
		t.Fatalf("fees lookup miss: %+v", got)
	}
	if _, ok := idx["nope"]; ok {
		t.Fatalf("unexpected hit for missing key")
	}
}
