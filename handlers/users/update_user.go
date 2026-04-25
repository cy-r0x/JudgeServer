package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	if userId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var payload UpdateUserPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	updates := make(map[string]interface{})
	if payload.Name != nil {
		updates["name"] = *payload.Name
	}
	if payload.Password != nil && *payload.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*payload.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to hash password", nil)
			return
		}
		updates["password"] = string(hashedPassword)
	}
	if payload.AdditionalInfo != nil {
		updates["additional_info"] = *payload.AdditionalInfo
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
		utils.SendResponse(w, http.StatusBadRequest, "No fields to update", nil)
		return
	}

	result := h.db.Model(&models.User{}).Where("id = ?", userId).Updates(updates)
	if result.Error != nil {
		log.Println(result.Error)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update user", nil)
		return
	}

	if result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusNotFound, "User not found", nil)
		return
	}

	utils.SendResponse(w, http.StatusOK, "User updated successfully", nil)
}
