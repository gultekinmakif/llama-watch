package dimensions

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
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
			path: writeJSON(t, root, "basic.json", map[string][]RawProtocol{
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
			path: writeJSON(t, root, "multi.json", map[string][]RawProtocol{
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
			path: writeJSON(t, root, "empty.json", map[string][]RawProtocol{}),
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
			path:    writeRaw(t, root, "bad.json", []byte("not json")),
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

// writeJSON writes a JSON-marshalled value to root/name and returns the absolute path.
func writeJSON(t *testing.T, root, name string, v any) string {
	t.Helper()
	full := filepath.Join(root, name)
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal %s: %v", name, err)
	}
	if err := os.WriteFile(full, b, 0o644); err != nil {
		t.Fatalf("write %s: %v", full, err)
	}
	return full
}

// writeRaw writes raw bytes to root/name and returns the absolute path.
func writeRaw(t *testing.T, root, name string, body []byte) string {
	t.Helper()
	full := filepath.Join(root, name)
	if err := os.WriteFile(full, body, 0o644); err != nil {
		t.Fatalf("write %s: %v", full, err)
	}
	return full
}
