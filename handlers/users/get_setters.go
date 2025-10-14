package users

import (
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetSetters(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, full_name, username FROM users WHERE role=$1`

	var users []UserResponse
	err := h.db.Select(&users, query, "setter")
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch setters")
		return
	}

	utils.SendResponse(w, http.StatusOK, users)
}
