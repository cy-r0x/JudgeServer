package cmd

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func Serve(PORT string) {
	//Init new Middleware Manager with Default Middlewares
	manager := middlewares.NewManager()
	manager.Use(middlewares.Prefilght, middlewares.Cors, middlewares.Logger)

	//Init New Mux and Init Routes
	mux := http.NewServeMux()
	initRoutes(mux, manager)

	//This will wrap the mux with global middlewares
	wrapedMux := manager.WrapMux(mux)

	http.ListenAndServe(PORT, wrapedMux)
}
