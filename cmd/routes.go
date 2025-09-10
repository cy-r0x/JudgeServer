package cmd

import (
	"net/http"

	"github.com/judgenot0/judge-backend/handlers"
	"github.com/judgenot0/judge-backend/middlewares"
)

func initRoutes(mux *http.ServeMux) {
	mngr := middlewares.NewManager()
	mngr.Use(middlewares.Logger)

	mux.Handle("GET /", mngr.With(handlers.HandleContests))

	mux.Handle("GET /contest/{contestId}", mngr.With(handlers.HandleProblemList))

	mux.Handle("GET /contest/{contestId}/{problemId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Fetch problem data from DB;
	}))

	mux.Handle("GET /setter-panel", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Auth User and return data from db else return 403;
	}))

	mux.Handle("GET /edit/{problemId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Auth the user->check if the user have access to that problem->fetch data-> return data
	}))

	mux.Handle("POST /submit/{problemId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Auth the user->check if the user have access to that problem->add to db->add to queue->return submission Id
	}))

	//TODO: More Routes to go
}
