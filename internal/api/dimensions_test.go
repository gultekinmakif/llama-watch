package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gultekinmakif/llama-watch/internal/models"
	"github.com/gultekinmakif/llama-watch/internal/registry"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func TestDimensionsShape(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/dimensions", nil)

	Dimensions(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type: want application/json; charset=utf-8, got %q", ct)
	}

	var resp dimensionsResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	cols := registry.Columns()
	if len(resp.Dimensions) != len(cols) {
		t.Fatalf("dimensions: want %d, got %d", len(cols), len(resp.Dimensions))
	}
	for i, d := range resp.Dimensions {
		if d.Kind == "" {
			t.Errorf("dimensions[%d].Kind: want non-empty, got empty", i)
		}
		if d.DisplayName == "" {
			t.Errorf("dimensions[%d].DisplayName: want non-empty, got empty", i)
		}
		if d.Coverage != 0 {
			t.Errorf("dimensions[%d].Coverage: want 0 (empty db), got %d", i, d.Coverage)
		}
	}
}

func TestDimensionsOrderMatchesRegistry(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/dimensions", nil)

	Dimensions(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, rec.Code)
	}

	var resp dimensionsResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	cols := registry.Columns()
	if len(resp.Dimensions) != len(cols) {
		t.Fatalf("dimensions length: want %d, got %d", len(cols), len(resp.Dimensions))
	}
	for i, c := range cols {
		if resp.Dimensions[i].Kind != c.Key {
			t.Errorf("dimensions[%d].Kind: want %q, got %q", i, c.Key, resp.Dimensions[i].Kind)
		}
		if resp.Dimensions[i].DisplayName != c.Label {
			t.Errorf("dimensions[%d].DisplayName: want %q, got %q", i, c.Label, resp.Dimensions[i].DisplayName)
		}
	}
}

func TestDimensionCoverageWithDB(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	withTx(t, func(tx *gorm.DB) {
		seedIdentity(t, tx, "p1", "P1", pq.StringArray{"ethereum"})
		seedIdentity(t, tx, "p2", "P2", pq.StringArray{"polygon"})

		for _, m := range []models.Matrix{
			{Slug: "p1", Metric: "tvl"},
			{Slug: "p1", Metric: "dailyFees"},
			{Slug: "p2", Metric: "tvl"},
		} {
			if err := tx.Create(&m).Error; err != nil {
				t.Fatalf("seed matrix %s/%s: %v", m.Slug, m.Metric, err)
			}
		}

		got, err := dimensionCoverage(t.Context(), tx)
		if err != nil {
			t.Fatalf("dimensionCoverage: %v", err)
		}
		if got["tvl"] != 2 {
			t.Errorf("tvl coverage: want 2, got %d", got["tvl"])
		}
		if got["dailyFees"] != 1 {
			t.Errorf("dailyFees coverage: want 1, got %d", got["dailyFees"])
		}
		if _, present := got["dailyRevenue"]; present {
			t.Errorf("dailyRevenue: want absent, got %d", got["dailyRevenue"])
		}
	})
}
