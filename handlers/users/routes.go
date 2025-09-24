package users

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle("POST /login", manager.With(h.Login))
	mux.Handle("POST /register", manager.With(h.Register))
	mux.Handle("POST /logout", manager.With(h.Logout))
	//TODO: More Routes to go
}
