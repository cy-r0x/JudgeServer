package users

import (
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	contestId := r.PathValue("contestId")
	if contestId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid contest ID", nil)
		return
	}

	var users []models.User
	// Fetch only specific fields or ignore password. With GORM, finding into User struct works naturally.
	err := h.db.Select("id", "name", "username", "room_no", "pc_no", "additional_info").Where("allowed_contest = ?", contestId).Find(&users).Error
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch users", nil)
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

	utils.SendResponse(w, http.StatusOK, "Users fetched successfully", response)
}
