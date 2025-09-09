package middlewares

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {

	controller := func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path, r.Method)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(controller)
}
