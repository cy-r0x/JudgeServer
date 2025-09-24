package contest

import (
	"net/http"

	"github.com/judgenot0/judge-backend/handlers/problem"
	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle("GET /contest/{contestId}", manager.With(problem.HandleProblemList))

	mux.Handle("GET /contest/{contestId}/{problemId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Fetch problem data from DB;
	}))
}
