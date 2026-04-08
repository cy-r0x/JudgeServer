package users

import (
	"log"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")

	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	result := h.db.Select("id", "full_name", "username", "role", "clan", "room_no", "pc_no", "allowed_contest", "created_at").Where("id = ?", userIdInt).First(&user)
	if result.Error != nil {
		log.Println(result.Error)
		utils.SendResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.SendResponse(w, http.StatusOK, user)
}
