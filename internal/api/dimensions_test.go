package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gultekinmakif/llama-watch/internal/registry"
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
