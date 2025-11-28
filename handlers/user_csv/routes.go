package usercsv

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /api/users/csv", manager.With(h.GetCSV, middlewares.Authenticate, middlewares.AuthenticateAdmin))
	mux.Handle("POST /api/users/register/csv", manager.With(h.AddUserCsv, middlewares.Authenticate, middlewares.AuthenticateAdmin))
}
