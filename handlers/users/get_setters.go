package users

import (
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetSetters(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	err := h.db.Select("id", "name", "username", "role", "additional_info").Where("role = ?", "setter").Find(&users).Error
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch setters", nil)
		return
	}

	var response []UserResponse
	for _, u := range users {
		response = append(response, UserResponse{
			Id:             u.Id,
			Name:           u.Name,
			Username:       u.Username,
			AdditionalInfo: u.AdditionalInfo,
		})
	}

	if response == nil {
		response = []UserResponse{}
	}

	utils.SendResponse(w, http.StatusOK, "Setters retrieved successfully", response)
}
