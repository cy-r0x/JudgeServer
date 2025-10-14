package users

import (
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Since we're using stateless JWT tokens, logout is handled client-side
	// by removing the token from storage. This endpoint acknowledges the logout.
	// TODO: Implement token blacklist if needed for enhanced security
	utils.SendResponse(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
