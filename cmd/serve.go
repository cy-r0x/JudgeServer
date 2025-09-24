package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/judgenot0/judge-backend/config"
	"github.com/judgenot0/judge-backend/handlers/contest"
	"github.com/judgenot0/judge-backend/handlers/problem"
	"github.com/judgenot0/judge-backend/handlers/setter"
	"github.com/judgenot0/judge-backend/handlers/submissions"
	"github.com/judgenot0/judge-backend/handlers/users"
	"github.com/judgenot0/judge-backend/middlewares"
)

func Serve() {
	config, err := config.GetConfig()
	if err != nil {
		os.Exit(1)
	}
	//Init new Middleware Manager with Default Middlewares
	manager := middlewares.NewManager(config)
	manager.Use(middlewares.Prefilght, middlewares.Cors, middlewares.Logger)

	contestHandler := contest.NewHandler()
	problemHandler := problem.NewHandler()
	setterHandler := setter.NewHandler()
	submissionsHandler := submissions.NewHandler()
	usersHandler := users.NewHandler(config)

	//Init New Mux and Init Routes
	mux := http.NewServeMux()
	contestHandler.RegisterRoutes(mux, manager)
	problemHandler.RegisterRoutes(mux, manager)
	setterHandler.RegisterRoutes(mux, manager)
	submissionsHandler.RegisterRoutes(mux, manager)
	usersHandler.RegisterRoutes(mux, manager)

	//This will wrap the mux with global middlewares
	wrapedMux := manager.WrapMux(mux)
	log.Printf("Server Running at http://localhost%s\n", config.HttpPort)
	http.ListenAndServe(config.HttpPort, wrapedMux)
}
