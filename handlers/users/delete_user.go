package users

import (
	"log"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")

	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	query := `DELETE FROM users WHERE id = $1`

	result, err := h.db.Exec(query, userIdInt)
	if err != nil {
		log.Println(err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to verify deletion")
		return
	}

	if rowsAffected == 0 {
		utils.SendResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.SendResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
