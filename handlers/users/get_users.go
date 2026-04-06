package users

import (
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	contestId := r.PathValue("contestId")

	// Convert contestId to int64
	contestIdInt, err := strconv.ParseInt(contestId, 10, 64)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid contest ID")
		return
	}

	var users []models.User
	// Fetch only specific fields or ignore password. With GORM, finding into User struct works naturally.
	err = h.db.Select("id", "full_name", "username", "clan", "room_no", "pc_no").Where("allowed_contest = ?", contestIdInt).Find(&users).Error
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	var response []User
	for _, u := range users {
		response = append(response, User{
			Id:       int64(u.ID),
			FullName: u.FullName,
			Username: u.Username,
			Clan:     u.Clan,
			RoomNo:   u.RoomNo,
			PcNo:     u.PcNo,
		})
	}

	if response == nil {
		response = []User{}
	}

	utils.SendResponse(w, http.StatusOK, response)
}
