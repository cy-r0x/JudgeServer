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
		utils.SendResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	result := h.db.Delete(&models.User{}, "id = ?", userId)
	if result.Error != nil {
		log.Println(result.Error)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete user", nil)
		return
	}

	if result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}

	utils.SendResponse(w, http.StatusOK, "User deleted successfully", nil)
}
