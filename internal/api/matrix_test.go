package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gultekinmakif/llama-watch/internal/registry"
)

func TestMatrixShape(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/matrix", nil)

	Matrix(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type: want application/json; charset=utf-8, got %q", ct)
	}

	var resp MatrixResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(resp.Columns) != len(registry.Columns()) {
		t.Fatalf("columns: want %d, got %d", len(registry.Columns()), len(resp.Columns))
	}
	if resp.Rows == nil {
		t.Error("rows: want non-nil empty slice, got nil")
	}
	if len(resp.Rows) != 0 {
		t.Errorf("rows: want empty, got len %d", len(resp.Rows))
	}
	if resp.Total != 0 {
		t.Errorf("total: want 0, got %d", resp.Total)
	}
}

func TestMatrixBadQuery(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/matrix?limit=2000", nil)

	Matrix(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var got struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Error.Code != "bad_request" {
		t.Errorf("code: want %q, got %q", "bad_request", got.Error.Code)
	}
	if got.Error.Message == "" {
		t.Error("message: want non-empty, got empty")
	}
}

func TestMatrixColumnsStableAcrossCalls(t *testing.T) {
	call := func() []registry.Column {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/matrix", nil)
		Matrix(rec, req)
		var resp MatrixResponse
		if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		return resp.Columns
	}

	first := call()
	second := call()

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("columns drifted between calls: first=%v second=%v", first, second)
	}
}
