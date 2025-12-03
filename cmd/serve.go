package cmd

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/judgenot0/judge-backend/config"
	"github.com/judgenot0/judge-backend/handlers/cluster"
	compilerun "github.com/judgenot0/judge-backend/handlers/compile_run"
	"github.com/judgenot0/judge-backend/handlers/contest"
	"github.com/judgenot0/judge-backend/handlers/contest_problems"
	"github.com/judgenot0/judge-backend/handlers/problem"
	"github.com/judgenot0/judge-backend/handlers/setter"
	"github.com/judgenot0/judge-backend/handlers/standings"
	"github.com/judgenot0/judge-backend/handlers/submissions"
	usercsv "github.com/judgenot0/judge-backend/handlers/user_csv"
	"github.com/judgenot0/judge-backend/handlers/users"
	"github.com/judgenot0/judge-backend/infra/db"
	"github.com/judgenot0/judge-backend/middlewares"
)

func Serve() {
	config, err := config.GetConfig()
	if err != nil {
		os.Exit(1)
	}

	dbConn, err := db.NewConnection(config.DB)
	if err != nil {
		os.Exit(1)
	}

	err = db.Migrate(dbConn, "./schema")
	if err != nil {
		os.Exit(1)
	}

	//Init new Middleware Manager with Default Middlewares
	manager := middlewares.NewManager()
	middlewares := middlewares.NewMiddlewares(config)

	manager.Use(middlewares.Cors, middlewares.Prefilght, middlewares.Logger)

	cluserHandler := cluster.NewHandler()
	contestHandler := contest.NewHandler(dbConn, config)
	contestProblemHandler := contest_problems.NewHandler(dbConn)
	problemHandler := problem.NewHandler(dbConn)
	setterHandler := setter.NewHandler(dbConn)
	standingsHandler := standings.NewHandler(dbConn)
	submissionsHandler := submissions.NewHandler(dbConn, config)
	usersHandler := users.NewHandler(config, dbConn)
	userCsvHandler := usercsv.NewHandler(dbConn)
	compilerunHandler := compilerun.NewHandler(dbConn, config)

	//Init New Mux and Init Routes
	mux := http.NewServeMux()
	cluserHandler.RegisterRoutes(mux, manager, middlewares)
	contestHandler.RegisterRoutes(mux, manager, middlewares)
	contestProblemHandler.RegisterRoute(mux, manager, middlewares)
	problemHandler.RegisterRoutes(mux, manager, middlewares)
	setterHandler.RegisterRoutes(mux, manager, middlewares)
	submissionsHandler.RegisterRoutes(mux, manager, middlewares)
	standingsHandler.RegisterRoutes(mux, manager, middlewares)
	usersHandler.RegisterRoutes(mux, manager, middlewares)
	userCsvHandler.RegisterRoutes(mux, manager, middlewares)
	compilerunHandler.RegisterRoute(mux, manager, middlewares)

	go func() {
		for {
			standingsHandler.MemoryEviction()
			time.Sleep(1 * time.Hour)
		}
	}()

	//This will wrap the mux with global middlewares
	wrapedMux := manager.WrapMux(mux)
	log.Printf("Server Running at http://localhost:%s\n", config.HttpPort)
	if err := http.ListenAndServe("0.0.0.0:"+config.HttpPort, wrapedMux); err != nil {
		log.Println("HTTP server error:", err)
		os.Exit(1)
	}
}
