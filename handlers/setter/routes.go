package setter

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager) {

	mux.Handle("GET /setter-panel", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Auth User and return data from db else return 403;
	}))

	mux.Handle("GET /edit/{problemId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Auth the user->check if the user have access to that problem->fetch data-> return data
	}))
}
