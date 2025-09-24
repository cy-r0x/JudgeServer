package submissions

import (
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux, manager *middlewares.Manager) {
	mux.Handle("GET /submissons", manager.With(h.GetSubmissions))
	mux.Handle("GET /submissions/{submissonId}", manager.With(h.GetSubmissionData))

	mux.Handle("POST /submissions/submit/{problemId}", manager.With(h.Submit, middlewares.Authenticate))
	//TODO: More Routes to go
}
