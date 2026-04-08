package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserPayload struct {
	FullName       *string `json:"full_name,omitempty"`
	Password       *string `json:"password,omitempty"`
	Clan           *string `json:"clan,omitempty"`
	RoomNo         *string `json:"room_no,omitempty"`
	PcNo           *string `json:"pc_no,omitempty"`
	AllowedContest *int64  `json:"allowed_contest,omitempty"`
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")

	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var payload UpdateUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updates := make(map[string]interface{})
	if payload.FullName != nil {
		updates["full_name"] = *payload.FullName
	}
	if payload.Password != nil && *payload.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*payload.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}
		updates["password"] = string(hashedPassword)
	}
	if payload.Clan != nil {
		updates["clan"] = *payload.Clan
	}
	if payload.RoomNo != nil {
		updates["room_no"] = *payload.RoomNo
	}
	if payload.PcNo != nil {
		updates["pc_no"] = *payload.PcNo
	}
	if payload.AllowedContest != nil {
		updates["allowed_contest"] = *payload.AllowedContest
	}

	if len(updates) == 0 {
		utils.SendResponse(w, http.StatusBadRequest, "No fields to update")
		return
	}

	result := h.db.Model(&models.User{}).Where("id = ?", userIdInt).Updates(updates)
	if result.Error != nil {
		log.Println(result.Error)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	if result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.SendResponse(w, http.StatusOK, "User updated successfully")
}
