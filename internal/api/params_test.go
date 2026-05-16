package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func parseURL(t *testing.T, rawurl string) MatrixQuery {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, rawurl, nil)
	got, err := ParseMatrixQuery(req)
	if err != nil {
		t.Fatalf("ParseMatrixQuery(%q): unexpected error: %v", rawurl, err)
	}
	return got
}

func TestParseMatrixQuery_HappyPaths(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want MatrixQuery
	}{
		{
			name: "empty defaults",
			url:  "/api/matrix",
			want: MatrixQuery{Limit: 200, Offset: 0},
		},
		{
			name: "limit and offset honored",
			url:  "/api/matrix?limit=500&offset=200",
			want: MatrixQuery{Limit: 500, Offset: 200},
		},
		{
			name: "limit boundary 1000 inclusive",
			url:  "/api/matrix?limit=1000",
			want: MatrixQuery{Limit: 1000},
		},
		{
			name: "limit boundary 1 inclusive",
			url:  "/api/matrix?limit=1",
			want: MatrixQuery{Limit: 1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseURL(t, tc.url)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("ParseMatrixQuery(%q):\n  got:  %+v\n  want: %+v", tc.url, got, tc.want)
			}
		})
	}
}

func TestParseMatrixQuery_BadPaths(t *testing.T) {
	cases := []struct {
		name        string
		url         string
		msgContains string
	}{
		{
			name:        "limit not integer",
			url:         "/api/matrix?limit=not_a_number",
			msgContains: "must be an integer",
		},
		{
			name:        "limit zero rejected",
			url:         "/api/matrix?limit=0",
			msgContains: "at least 1",
		},
		{
			name:        "limit over max rejected",
			url:         "/api/matrix?limit=1001",
			msgContains: "at most 1000",
		},
		{
			name:        "offset negative rejected",
			url:         "/api/matrix?offset=-5",
			msgContains: "at least 0",
		},
		{
			name:        "offset not integer",
			url:         "/api/matrix?offset=not_a_number",
			msgContains: "must be an integer",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			_, err := ParseMatrixQuery(req)
			if err == nil {
				t.Fatalf("ParseMatrixQuery(%q): want error, got nil", tc.url)
			}
			var perr *ParseError
			if !errors.As(err, &perr) {
				t.Fatalf("ParseMatrixQuery(%q): want *ParseError, got %T", tc.url, err)
			}
			if perr.Code != "bad_request" {
				t.Errorf("code: want %q, got %q", "bad_request", perr.Code)
			}
			if !strings.Contains(perr.Message, tc.msgContains) {
				t.Errorf("message: want substring %q, got %q", tc.msgContains, perr.Message)
			}
		})
	}
}
