package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lib/pq"
	"gorm.io/gorm"
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

func TestListChainsWithDB(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	withTx(t, func(tx *gorm.DB) {
		seedIdentity(t, tx, "p1", "P1", pq.StringArray{"ethereum", "polygon"})
		seedIdentity(t, tx, "p2", "P2", pq.StringArray{"ethereum"})
		seedIdentity(t, tx, "p3", "P3", pq.StringArray{"bsc"})

		rows, err := listChains(t.Context(), tx)
		if err != nil {
			t.Fatalf("listChains: %v", err)
		}

		got := make(map[string]int, len(rows))
		for _, r := range rows {
			got[r.Key] = r.ProtocolCount
		}
		if got["ethereum"] != 2 {
			t.Errorf("ethereum count: want 2, got %d", got["ethereum"])
		}
		if got["polygon"] != 1 {
			t.Errorf("polygon count: want 1, got %d", got["polygon"])
		}
		if got["bsc"] != 1 {
			t.Errorf("bsc count: want 1, got %d", got["bsc"])
		}

		for _, r := range rows {
			if r.Key == "ethereum" && r.Label != "Ethereum" {
				t.Errorf("ethereum label: want Ethereum, got %s", r.Label)
			}
		}
	})
}
