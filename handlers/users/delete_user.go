package users

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	if userId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	result := h.db.Delete(&models.User{}, "id = ?", userId)
	if result.Error != nil {
		log.Println(result.Error)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	if result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.SendResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
