package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChainsShape(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/chains", nil)

	Chains(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type: want application/json; charset=utf-8, got %q", ct)
	}

	var resp chainsResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Chains == nil {
		t.Error("chains: want non-nil empty slice, got nil")
	}
}

func TestChainsLabelTitlecase(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"ethereum", "Ethereum"},
		{"bsc", "Bsc"},
		{"", ""},
		{"a", "A"},
	}
	for _, c := range cases {
		if got := titlecase(c.in); got != c.want {
			t.Errorf("titlecase(%q): want %q, got %q", c.in, c.want, got)
		}
	}
}
