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
		{"multi-version sibling v4", "Uniswap V4", "uniswap-v4"},
		{"apostrophe stripped", "d'CENT", "dcent"},
		{"multiple apostrophes stripped", "L'Anza's Vault", "lanzas-vault"},
		{"multiple spaces each become dash", "Lido Finance V2", "lido-finance-v2"},
		{"adjacent spaces produce adjacent dashes (upstream parity)", "Aave  V2", "aave--v2"},
		{"empty", "", ""},
		{"already-canonical untouched", "uniswap-v2", "uniswap-v2"},
		{"punctuation other than apostrophe survives", "0x_Protocol", "0x_protocol"},
		{"unicode lowercases via Go case folding", "Café", "café"},
		{"all-whitespace becomes all-dashes", "   ", "---"},
		{"tab survives (only ASCII space is replaced)", "foo\tbar", "foo\tbar"},
		{"newline survives (only ASCII space is replaced)", "foo\nbar", "foo\nbar"},
	}
	for _, tc := range cases {
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
		{"underscore filename with ts ext", "Uniswap_V2.ts", "uniswap-v2"},
		{"collapse adjacent dashes", "uniswap--v2", "uniswap-v2"},
		{"trim leading and trailing dashes", "-uniswap-v2-", "uniswap-v2"},
		{"underscore plus collapse", "uniswap__v2.ts", "uniswap-v2"},
		{"mixed underscore-dash collapse", "foo_-bar.ts", "foo-bar"},
		{"space treated as dash", "Aave V2.ts", "aave-v2"},
		{"only an extension", ".ts", ""},
		{"empty", "", ""},
		{"no extension lowercased", "WBTC", "wbtc"},
		{"apostrophe preserved (diverges from Canonical)", "d'cent.ts", "d'cent"},
		{"multi-dot extension only strips one suffix", "foo.test.ts", "foo.test"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FromFilename(tc.in)
			if got != tc.want {
				t.Fatalf("FromFilename(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
