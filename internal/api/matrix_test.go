package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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

	if len(resp.Columns) != 9 {
		t.Fatalf("columns: want 9, got %d", len(resp.Columns))
	}
	if resp.Columns[0].Key != "tvl" {
		t.Errorf("columns[0].key: want %q, got %q", "tvl", resp.Columns[0].Key)
	}
	if resp.Columns[8].Key != "dailyActiveUsers" {
		t.Errorf("columns[8].key: want %q, got %q", "dailyActiveUsers", resp.Columns[8].Key)
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

func TestMatrixColumnsStableAcrossCalls(t *testing.T) {
	call := func() []Column {
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
