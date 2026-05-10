package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

func (m *Middlewares) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		if r.URL.Path != "/get_nodes" {
			slog.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start).String(),
			)
		}
	})
}
