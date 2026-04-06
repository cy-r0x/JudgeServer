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
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	contest.StartTime = contest.StartTime.UTC()

	var description *string
	if contest.Description != "" {
		description = &contest.Description
	}

	updateData := models.Contest{
		Title:           contest.Title,
		Description:     description,
		StartTime:       contest.StartTime,
		DurationSeconds: contest.DurationSeconds,
	}

	result := h.db.Model(&models.Contest{}).Where("id = ?", contest.Id).Updates(updateData)
	if result.Error != nil || result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update contest")
		return
	}

	var updatedContest models.Contest
	h.db.First(&updatedContest, contest.Id)

	contest.CreatedAt = updatedContest.CreatedAt

	utils.SendResponse(w, http.StatusOK, contest)
}
