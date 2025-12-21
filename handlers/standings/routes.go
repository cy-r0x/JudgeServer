package standings

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /api/standings/{contestId}", manager.With(h.GetStandings))
	mux.Handle("GET /api/standings/export/{contestId}", manager.With(h.ExportStandings, middlewares.Authenticate, middlewares.AuthenticateAdmin))
}
