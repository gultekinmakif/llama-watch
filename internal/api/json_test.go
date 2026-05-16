package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSONSuccess(t *testing.T) {
	rec := httptest.NewRecorder()
	writeJSON(rec, http.StatusOK, map[string]string{"hello": "world"})

	if rec.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type: want application/json; charset=utf-8, got %q", ct)
	}
	if got, want := rec.Body.String(), "{\"hello\":\"world\"}\n"; got != want {
		t.Fatalf("body: want %q, got %q", want, got)
	}
}

func TestWriteJSONEncodeFailure(t *testing.T) {
	rec := httptest.NewRecorder()
	// channels are unsupported by encoding/json and force Encode to fail.
	writeJSON(rec, http.StatusOK, make(chan int))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: want %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type: want application/json; charset=utf-8, got %q", ct)
	}
	want := `{"error":{"code":"internal","message":"encode failed"}}` + "\n"
	if got := rec.Body.String(); got != want {
		t.Fatalf("body: want %q, got %q", want, got)
	}
}
