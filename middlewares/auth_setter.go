package middlewares

import (
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (m *Middlewares) AuthenticateSetter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, ok := r.Context().Value("user").(*Payload)
		if (!ok) {
			utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
			return
		}
		if payload.Role != "setter" {
			utils.SendResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}
