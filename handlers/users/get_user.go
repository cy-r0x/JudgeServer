package users

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	if userId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var user models.User
	result := h.db.Select("id", "name", "username", "role", "additional_info", "room_no", "pc_no", "allowed_contest", "created_at").Where("id = ?", userId).First(&user)
	if result.Error != nil {
		log.Println(result.Error)
		utils.SendResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}

	utils.SendResponse(w, http.StatusOK, "User fetched successfully", user)
}
