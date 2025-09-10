package cmd

import (
	"net/http"
)

func Serve() {
	mux := http.NewServeMux()
	initRoutes(mux)
	http.ListenAndServe(":8080", mux)
}
