package submissions

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /submissions", manager.With(h.GetSubmissions, middlewares.Authenticate))
	mux.Handle("GET /submissions/{submissonId}", manager.With(h.GetSubmissionData, middlewares.Authenticate))

	mux.Handle("POST /submissions/submit/{problemId}", manager.With(h.Submit, middlewares.Authenticate))
	//TODO: More Routes to go
}
