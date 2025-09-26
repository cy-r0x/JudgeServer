package contest

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /contests", manager.With(h.ListContests))
	mux.Handle("GET /contests/{contestId}", manager.With(h.GetContest))
	mux.Handle("POST /contests/create", manager.With(h.CreateContest, middlewares.Authenticate))
	mux.Handle("POST /contests/update", manager.With(h.UpdateContest, middlewares.Authenticate))
	mux.Handle("DELETE /contests/delete", manager.With(h.DeleteContest, middlewares.Authenticate))
}
