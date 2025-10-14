package problem

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /api/problems/{problemId}", manager.With(h.GetProblem, middlewares.Authenticate))
	mux.Handle("POST /api/problems", manager.With(h.CreateProblem, middlewares.Authenticate, middlewares.AuthenticateSetter))
	mux.Handle("PUT /api/problems", manager.With(h.UpdateProblem, middlewares.Authenticate, middlewares.AuthenticateSetter))
}
