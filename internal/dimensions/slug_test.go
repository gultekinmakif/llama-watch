package dimensions

import "testing"

func TestCanonical(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"lowercases", "Uniswap", "uniswap"},
		{"multi-version sibling v2", "Uniswap V2", "uniswap-v2"},
		{"multi-version sibling v3", "Uniswap V3", "uniswap-v3"},
		{"apostrophe stripped", "d'CENT", "dcent"},
		{"multiple spaces each become dash", "Lido Finance V2", "lido-finance-v2"},
		{"already-canonical untouched", "uniswap-v2", "uniswap-v2"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := Canonical(tc.in)
			if got != tc.want {
				t.Fatalf("Canonical(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestFromFilename(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain ts basename", "uniswap-v2.ts", "uniswap-v2"},
		{"plain js basename", "wbtc.js", "wbtc"},
		{"underscore to dash", "Uniswap_V2", "uniswap-v2"},
		{"space treated as dash", "Aave V2.ts", "aave-v2"},
		{"no extension lowercased", "WBTC", "wbtc"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := FromFilename(tc.in)
			if got != tc.want {
				t.Fatalf("FromFilename(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
