package cmd

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func Serve(HTTP_PORT string) {
	//Init new Middleware Manager with Default Middlewares
	manager := middlewares.NewManager()
	manager.Use(middlewares.Prefilght, middlewares.Cors, middlewares.Logger)

	//Init New Mux and Init Routes
	mux := http.NewServeMux()
	initRoutes(mux, manager)

	//This will wrap the mux with global middlewares
	wrapedMux := manager.WrapMux(mux)
	log.Printf("Server Running at http://localhost%s", HTTP_PORT)
	http.ListenAndServe(HTTP_PORT, wrapedMux)
}
