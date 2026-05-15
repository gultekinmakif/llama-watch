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
