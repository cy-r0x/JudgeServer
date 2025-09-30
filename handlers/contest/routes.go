package contest

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /api/contests", manager.With(h.ListContests))
	mux.Handle("GET /api/contests/{contestId}", manager.With(h.GetContest))
	mux.Handle("POST /api/contests/create", manager.With(h.CreateContest, middlewares.Authenticate, middlewares.AuthenticateAdmin))
	mux.Handle("POST /api/contests/update", manager.With(h.UpdateContestIndex, middlewares.Authenticate, middlewares.AuthenticateAdmin))
}
