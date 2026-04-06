package contest_problems

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoute(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /api/contests/problems/{contestId}", manager.With(h.GetContestProblems, middlewares.Authenticate))
	mux.Handle("POST /api/contests/assign", manager.With(h.AssignContestProblems, middlewares.Authenticate, middlewares.AuthenticateAdmin))
	mux.Handle("PUT /api/contests/index", manager.With(h.UpdateContestIndex, middlewares.Authenticate, middlewares.AuthenticateAdmin))
}
