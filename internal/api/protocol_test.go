package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProtocolNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/protocols/nope", nil)
	req.SetPathValue("slug", "nope")

	Protocol(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type: want application/json; charset=utf-8, got %q", ct)
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
	if got.Error.Code != "not_found" {
		t.Errorf("code: want %q, got %q", "not_found", got.Error.Code)
	}
	if got.Error.Message == "" {
		t.Error("message: want non-empty, got empty")
	}
}

