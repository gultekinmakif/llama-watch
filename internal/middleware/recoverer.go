// 2.3: Outermost middleware. Catches panics, returns 500, logs stack + request_id.
package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Body literal (not extracted to internal/api) to avoid an api→middleware→api import cycle.
var panicBody = []byte(`{"error":{"code":"internal","message":"internal server error"}}` + "\n")

// Recoverer recovers from panics, logs the panic with stack trace and request
// ID (if present), and returns a 500 to the client.
func Recoverer(next http.Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic",
					"error", err,
					"path", r.URL.Path,
					"request_id", GetReqID(r.Context()),
					"stack", string(debug.Stack()),
				)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write(panicBody)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
