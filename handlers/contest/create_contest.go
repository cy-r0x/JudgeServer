package contest

import (
	"encoding/json"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) CreateContest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reqContest Contest
	err := decoder.Decode(&reqContest)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Convert start time to UTC
	startTime := reqContest.StartTime.UTC()
	var description *string
	if reqContest.Description != "" {
		description = &reqContest.Description
	}

	newContest := models.Contest{
		Title:           reqContest.Title,
		Description:     description,
		StartTime:       startTime,
		DurationSeconds: reqContest.DurationSeconds,
	}

	err = h.db.Create(&newContest).Error
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create contest")
		return
	}

	reqContest.Id = newContest.ID
	reqContest.CreatedAt = newContest.CreatedAt
	reqContest.StartTime = newContest.StartTime

	utils.SendResponse(w, http.StatusCreated, reqContest)
}
