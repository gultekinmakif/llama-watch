package dimensions

import (
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
