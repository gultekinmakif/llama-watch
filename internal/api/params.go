// Query-param parser and validation for /api/matrix.
package api

import (
	"net/http"
	"strconv"
	"strings"
)

type MatrixQuery struct {
	Limit      int      // 1..1000
	Offset     int      // >= 0
	Chains     []string // lowercased, deduped, order-preserving; nil if absent
	Categories []string // case preserved, deduped, order-preserving; nil if absent
	Q          string   // trimmed; empty if absent
	Sort       string   // one of: name, category, coverage; never empty post-parse
	Order      string   // "asc" or "desc"; never empty post-parse
}

type ParseError struct {
	Code    string
	Message string
}

func (e *ParseError) Error() string { return e.Message }

// parseFail returns the zero MatrixQuery plus a bad_request *ParseError wrapping msg.
// All parser failures route through here so the code field and the zero-value
// invariant on the error path cannot drift.
func parseFail(msg string) (MatrixQuery, error) {
	return MatrixQuery{}, &ParseError{Code: "bad_request", Message: msg}
}

const (
	defaultLimit = 200
	maxLimit     = 1000
)

// Whitelist of sort keys. Referenced by both the parser and the error message
// so they cannot drift apart.
var sortWhitelist = []string{"name", "category", "coverage"}

// Default order per sort key when ?order= is absent.
var defaultOrderBySort = map[string]string{
	"name":     "asc",
	"category": "asc",
	"coverage": "desc",
}

// ParseMatrixQuery parses and validates the /api/matrix query string.
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
			return parseFail("invalid limit: must be an integer")
		}
		if n < 1 {
			return parseFail("invalid limit: must be at least 1")
		}
		if n > maxLimit {
			return parseFail("invalid limit: must be at most 1000")
		}
		out.Limit = n
	}

	offsetRaw := q.Get("offset")
	if offsetRaw == "" {
		out.Offset = 0
	} else {
		n, err := strconv.Atoi(offsetRaw)
		if err != nil {
			return parseFail("invalid offset: must be an integer")
		}
		if n < 0 {
			return parseFail("invalid offset: must be at least 0")
		}
		out.Offset = n
	}

	if raw, ok := q["chains"]; ok && len(raw) > 0 {
		out.Chains = parseCSV(raw[0], true)
	}

	if raw, ok := q["categories"]; ok && len(raw) > 0 {
		out.Categories = parseCSV(raw[0], false)
	}

	out.Q = strings.TrimSpace(q.Get("q"))

	sortRaw := q.Get("sort")
	if sortRaw == "" {
		out.Sort = "coverage"
	} else if _, ok := defaultOrderBySort[sortRaw]; ok {
		out.Sort = sortRaw
	} else {
		return parseFail("invalid sort: must be one of " + strings.Join(sortWhitelist, ", "))
	}

	orderRaw := q.Get("order")
	switch orderRaw {
	case "":
		out.Order = defaultOrderBySort[out.Sort]
	case "asc", "desc":
		out.Order = orderRaw
	default:
		return parseFail("invalid order: must be asc or desc")
	}

	return out, nil
}

// parseCSV splits a CSV value, drops empty entries, and dedupes preserving
// first occurrence. When lower is true, each entry is lowercased before dedup.
func parseCSV(raw string, lower bool) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if lower {
			p = strings.ToLower(p)
		}
		if p == "" {
			continue
		}
		if _, dup := seen[p]; dup {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
