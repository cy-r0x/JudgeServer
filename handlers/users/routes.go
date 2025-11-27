package users

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /api/users/{contestId}", manager.With(h.GetUsers, middlewares.Authenticate, middlewares.AuthenticateAdmin))
	mux.Handle("GET /api/users/setter", manager.With(h.GetSetters, middlewares.Authenticate, middlewares.AuthenticateAdmin))
	mux.Handle("GET /api/users/csv", manager.With(h.GetCSV, middlewares.Authenticate, middlewares.AuthenticateAdmin))
	mux.Handle("POST /api/users/login", manager.With(h.Login))
	mux.Handle("POST /api/users/register", manager.With(h.CreateUser, middlewares.Authenticate, middlewares.AuthenticateAdmin))
	mux.Handle("POST /api/users/register/csv", manager.With(h.AddUserCsv, middlewares.Authenticate, middlewares.AuthenticateAdmin))
	mux.Handle("POST /api/users/delete/{userId}", manager.With(h.DeleteUser, middlewares.Authenticate, middlewares.AuthenticateAdmin))
}
