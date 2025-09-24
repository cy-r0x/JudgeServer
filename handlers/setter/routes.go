package setter

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle("GET /setter-panel", manager.With(h.GetSetterData))
	mux.Handle("GET /edit/{problemId}", manager.With(h.EditProblem))
	mux.Handle("POST /edit/{problemId}", manager.With(h.SaveProblem))
	mux.Handle("DELETE /edit/{problemId}", manager.With(h.DeleteProblem))
}
