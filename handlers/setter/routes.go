package setter

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /setter-panel", manager.With(h.GetSetterData, middlewares.Authenticate))
	mux.Handle("GET /edit/{problemId}", manager.With(h.EditProblem, middlewares.Authenticate))
	mux.Handle("POST /edit/{problemId}", manager.With(h.SaveProblem, middlewares.Authenticate))
	mux.Handle("DELETE /edit/{problemId}", manager.With(h.DeleteProblem, middlewares.Authenticate))
}
