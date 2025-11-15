package cluster

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("POST /register_node", manager.With(h.registerNode))
	mux.Handle("GET /get_nodes", manager.With(h.getNodes))
}
