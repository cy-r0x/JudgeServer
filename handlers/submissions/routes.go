package submissions

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /submissions", manager.With(h.ListUserSubmissions, middlewares.Authenticate))
	mux.Handle("GET /submissions/{submissonId}", manager.With(h.GetSubmission, middlewares.Authenticate))
	mux.Handle("GET /submissions/all/{contestId}", manager.With(h.ListAllSubmissions, middlewares.Authenticate))
	mux.Handle("POST /submissions/submit/{problemId}", manager.With(h.CreateSubmission, middlewares.Authenticate))
}
