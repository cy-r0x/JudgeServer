package contest

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /contests", manager.With(h.GetContests))
	mux.Handle("GET /contests/{contestId}", manager.With(h.GetContestData))
}
