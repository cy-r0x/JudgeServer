package users

import (
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetSetters(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	err := h.db.Select("id", "full_name", "username", "role").Where("role = ?", "setter").Find(&users).Error
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch setters")
		return
	}

	var response []UserResponse
	for _, u := range users {
		response = append(response, UserResponse{
			Id:       u.ID,
			FullName: u.FullName,
			Username: u.Username,
			Clan:     u.Clan,
		})
	}

	if response == nil {
		response = []UserResponse{}
	}

	utils.SendResponse(w, http.StatusOK, response)
}
