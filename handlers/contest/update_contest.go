package contest

import (
	"encoding/json"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) UpdateContest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var contest Contest
	err := decoder.Decode(&contest)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	contest.StartTime = contest.StartTime.UTC()
	contest.EndTime = contest.EndTime.UTC()

	var description *string
	if contest.Description != "" {
		description = &contest.Description
	}

	updateData := map[string]interface{}{
		"title":            contest.Title,
		"user_prefix":      contest.UserPrefix,
		"description":      description,
		"start_time":       contest.StartTime,
		"end_time":         contest.EndTime,
		"duration_seconds": contest.DurationSeconds,
	}

	result := h.db.Model(&models.Contest{}).Where("id = ?", contest.Id).Updates(updateData)
	if result.Error != nil || result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update contest", nil)
		return
	}

	var updatedContest models.Contest
	if err := h.db.Where("id = ?", contest.Id).First(&updatedContest).Error; err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch updated contest", nil)
		return
	}

	contest.CreatedAt = updatedContest.CreatedAt

	utils.SendResponse(w, http.StatusOK, nil, contest)
}
