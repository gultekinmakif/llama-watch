package dimensions

import (
	"reflect"
	"testing"

	"github.com/gultekinmakif/llama-watch/internal/testutil"
)

func TestLoadProtocols(t *testing.T) {
	root := t.TempDir()

	cases := []struct {
		name    string
		path    string
		wantErr bool
		check   func(t *testing.T, got map[string][]RawProtocol)
	}{
		{
			name: "basic shape preserves keys and dimensions",
			path: testutil.WriteJSON(t, root, "basic.json", map[string][]RawProtocol{
				"data1": {{
					Name:       "Uniswap V2",
					Category:   "Dexes",
					Chains:     []string{"ethereum", "polygon"},
					Module:     "uniswap-v2/index.js",
					Dimensions: map[string]string{"fees": "uniswap-v2", "dexs": "uniswap-v2"},
				}},
				"data2": {{
					Name:     "WBTC",
					Category: "Bridge",
					Chains:   []string{"ethereum"},
					Module:   "wbtc.js",
				}},
			}),
			check: func(t *testing.T, got map[string][]RawProtocol) {
				if len(got) != 2 {
					t.Fatalf("want 2 source keys, got %d", len(got))
				}
				if got["data1"][0].Name != "Uniswap V2" {
					t.Fatalf("data1[0].Name = %q", got["data1"][0].Name)
				}
				if got["data2"][0].Name != "WBTC" {
					t.Fatalf("data2[0].Name = %q", got["data2"][0].Name)
				}
				wantDims := map[string]string{"fees": "uniswap-v2", "dexs": "uniswap-v2"}
				if !reflect.DeepEqual(got["data1"][0].Dimensions, wantDims) {
					t.Fatalf("dims = %v", got["data1"][0].Dimensions)
				}
			},
		},
		{
			name: "multi version siblings preserved",
			path: testutil.WriteJSON(t, root, "multi.json", map[string][]RawProtocol{
				"data1": {
					{Name: "Uniswap V2", Module: "uniswap-v2/index.js"},
					{Name: "Uniswap V3", Module: "uniswap-v3/index.js"},
					{Name: "Uniswap V4", Module: "uniswap-v4/index.js"},
				},
			}),
			check: func(t *testing.T, got map[string][]RawProtocol) {
				if len(got["data1"]) != 3 {
					t.Fatalf("want 3 siblings, got %d", len(got["data1"]))
				}
			},
		},
		{
			name: "empty json returns empty map",
			path: testutil.WriteJSON(t, root, "empty.json", map[string][]RawProtocol{}),
			check: func(t *testing.T, got map[string][]RawProtocol) {
				if len(got) != 0 {
					t.Fatalf("want 0, got %d", len(got))
				}
			},
		},
		{
			name:    "missing file errors",
			path:    "/nonexistent/protocols.json",
			wantErr: true,
		},
		{
			name:    "malformed json errors",
			path:    testutil.WriteFile(t, root, "bad.json", "not json"),
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := LoadProtocols(tc.path)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("LoadProtocols: %v", err)
			}
			tc.check(t, got)
		})
	}
}
