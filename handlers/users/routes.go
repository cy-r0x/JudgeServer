package users

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle("POST /login", manager.With(h.HandleLogin))

	mux.Handle("POST /submissions/submit/{problemId}", manager.With(h.HandleSubmit, middlewares.Authenticate))

	//TODO: More Routes to go
}
