package middlewares

import (
	"log"
	"net/http"
	"time"
)

func (m *Middlewares) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current_time := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.URL.Path, r.Method, time.Since(current_time))
	})
}
