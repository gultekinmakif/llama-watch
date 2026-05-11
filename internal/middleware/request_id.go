package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

type ctxKey string

// Key to use when setting the request ID.
const RequestIDKey ctxKey = "reqID"
const RequestIDHeader = "X-Request-ID"

var prefix string
var reqid uint64

/*
A quick note on the statistics here: we're trying to calculate the chance that
two randomly generated base62 prefixes will collide. We use the formula from
http://en.wikipedia.org/wiki/Birthday_problem

P[m, n] \approx 1 - e^{-m^2/2n}

We ballpark an upper bound for $m$ by imagining (for whatever reason) a server
that restarts every second over 10 years, for $m = 86400 * 365 * 10 = 315360000$

For a $k$ character base-62 identifier, we have $n(k) = 62^k$

Plugging this in, we find $P[m, n(10)] \approx 5.75%$, which is good enough for
our purposes, and is surely more than anyone would ever need in practice -- a
process that is rebooted a handful of times a day for a hundred years has less
than a millionth of a percent chance of generating two colliding IDs.
*/

func init() {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		_, _ = rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

// RequestID is a middleware that injects a request ID into the context of each
// request. If the incoming request carries an X-Request-ID header (e.g. from
// an upstream LB or gateway), it is preserved; otherwise a new ID is
// generated of the form "host.example.com/random-0001", where "random" is a
// base62 random string that uniquely identifies this go process, and where
// the last number is an atomically incremented request counter.
func RequestID(next http.Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if id == "" {
			id = fmt.Sprintf("%s-%06d", prefix, atomic.AddUint64(&reqid, 1))
		}
		w.Header().Set(RequestIDHeader, id)
		ctx := context.WithValue(r.Context(), RequestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func GetReqID(ctx context.Context) string {
	if v, ok := ctx.Value(RequestIDKey).(string); ok {
		return v
	}
	return ""
}
