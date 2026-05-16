// Query-param parser and validation for /api/matrix.
package api

import (
	"net/http"
	"strconv"
)

type MatrixQuery struct {
	Limit  int // 1..1000
	Offset int // >= 0
}

type ParseError struct {
	Code    string
	Message string
}

func (e *ParseError) Error() string { return e.Message }

const (
	defaultLimit = 200
	maxLimit     = 1000
)

// ParseMatrixQuery parses and validates ?limit and ?offset from r.URL.Query().
// Returns a typed MatrixQuery on success, or a *ParseError on validation failure.
func ParseMatrixQuery(r *http.Request) (MatrixQuery, error) {
	q := r.URL.Query()
	out := MatrixQuery{}

	limitRaw := q.Get("limit")
	if limitRaw == "" {
		out.Limit = defaultLimit
	} else {
		n, err := strconv.Atoi(limitRaw)
		if err != nil {
			return MatrixQuery{}, &ParseError{Code: "bad_request", Message: "invalid limit: must be an integer"}
		}
		if n < 1 {
			return MatrixQuery{}, &ParseError{Code: "bad_request", Message: "invalid limit: must be at least 1"}
		}
		if n > maxLimit {
			return MatrixQuery{}, &ParseError{Code: "bad_request", Message: "invalid limit: must be at most 1000"}
		}
		out.Limit = n
	}

	offsetRaw := q.Get("offset")
	if offsetRaw == "" {
		out.Offset = 0
	} else {
		n, err := strconv.Atoi(offsetRaw)
		if err != nil {
			return MatrixQuery{}, &ParseError{Code: "bad_request", Message: "invalid offset: must be an integer"}
		}
		if n < 0 {
			return MatrixQuery{}, &ParseError{Code: "bad_request", Message: "invalid offset: must be at least 0"}
		}
		out.Offset = n
	}

	return out, nil
}
