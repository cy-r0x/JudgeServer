package cmd

import (
	"net/http"

	"github.com/judgenot0/judge-backend/handlers"
	"github.com/judgenot0/judge-backend/middlewares"
)

func initRoutes(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle("GET /", manager.With(handlers.HandleContests))
	mux.Handle("POST /login", manager.With(handlers.HandleLogin))

	mux.Handle("GET /contest/{contestId}", manager.With(handlers.HandleProblemList))

	mux.Handle("GET /contest/{contestId}/{problemId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Fetch problem data from DB;
	}))

	mux.Handle("GET /setter-panel", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Auth User and return data from db else return 403;
	}))

	mux.Handle("GET /edit/{problemId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Auth the user->check if the user have access to that problem->fetch data-> return data
	}))

	mux.Handle("POST /submissions/submit/{problemId}", manager.With(handlers.HandleSubmit, middlewares.Authenticate))

	//TODO: More Routes to go
}
