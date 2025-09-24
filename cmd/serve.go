package cmd

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/handlers/contest"
	"github.com/judgenot0/judge-backend/handlers/problem"
	"github.com/judgenot0/judge-backend/handlers/root"
	"github.com/judgenot0/judge-backend/handlers/setter"
	"github.com/judgenot0/judge-backend/handlers/users"
	"github.com/judgenot0/judge-backend/middlewares"
)

func Serve(HTTP_PORT string) {
	//Init new Middleware Manager with Default Middlewares
	manager := middlewares.NewManager()
	manager.Use(middlewares.Prefilght, middlewares.Cors, middlewares.Logger)

	contestHandler := contest.NewHandler()
	problemHandler := problem.NewHandler()
	rootHandler := root.NewHandler()
	setterHandler := setter.NewHandler()
	usersHandler := users.NewHandler()

	//Init New Mux and Init Routes
	mux := http.NewServeMux()
	contestHandler.RegisterRoutes(mux, manager)
	problemHandler.RegisterRoutes(mux, manager)
	rootHandler.RegisterRoutes(mux, manager)
	setterHandler.RegisterRoutes(mux, manager)
	usersHandler.RegisterRoutes(mux, manager)

	//This will wrap the mux with global middlewares
	wrapedMux := manager.WrapMux(mux)
	log.Printf("Server Running at http://localhost%s\n", HTTP_PORT)
	http.ListenAndServe(HTTP_PORT, wrapedMux)
}
