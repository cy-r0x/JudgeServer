package compilerun

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoute(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("POST /api/compile", manager.With(h.CompileRun, middlewares.Authenticate))
}
