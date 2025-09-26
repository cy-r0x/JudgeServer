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
	"github.com/judgenot0/judge-backend/infra/db"
	"github.com/judgenot0/judge-backend/middlewares"
)

func Serve() {
	config, err := config.GetConfig()
	if err != nil {
		os.Exit(1)
	}

	dbConn, err := db.NewConnection()
	if err != nil {
		os.Exit(1)
	}

	//Init new Middleware Manager with Default Middlewares
	manager := middlewares.NewManager()
	middlewares := middlewares.NewMiddlewares(config)

	manager.Use(middlewares.Prefilght, middlewares.Cors, middlewares.Logger)

	contestHandler := contest.NewHandler(dbConn)
	problemHandler := problem.NewHandler(dbConn)
	setterHandler := setter.NewHandler(dbConn)
	submissionsHandler := submissions.NewHandler(dbConn)
	usersHandler := users.NewHandler(config, dbConn)

	//Init New Mux and Init Routes
	mux := http.NewServeMux()
	contestHandler.RegisterRoutes(mux, manager, middlewares)
	problemHandler.RegisterRoutes(mux, manager, middlewares)
	setterHandler.RegisterRoutes(mux, manager, middlewares)
	submissionsHandler.RegisterRoutes(mux, manager, middlewares)
	usersHandler.RegisterRoutes(mux, manager, middlewares)

	//This will wrap the mux with global middlewares
	wrapedMux := manager.WrapMux(mux)
	log.Printf("Server Running at http://localhost%s\n", config.HttpPort)
	http.ListenAndServe(config.HttpPort, wrapedMux)
}
