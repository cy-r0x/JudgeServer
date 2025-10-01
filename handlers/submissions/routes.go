package submissions

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager, middlewares *middlewares.Middlewares) {
	mux.Handle("GET /api/submissions", manager.With(h.ListUserSubmissions, middlewares.Authenticate))
	mux.Handle("POST /api/submissions", manager.With(h.CreateSubmission, middlewares.Authenticate))
	mux.Handle("PUT /api/submissions", manager.With(h.UpdateSubmission, middlewares.AuthEngine))
	mux.Handle("GET /api/submissions/{submissonId}", manager.With(h.GetSubmission, middlewares.Authenticate))
	mux.Handle("GET /api/submissions/all/{contestId}", manager.With(h.ListAllSubmissions, middlewares.Authenticate))
}
