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
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	startTime := reqContest.StartTime.UTC()
	endTime := reqContest.EndTime.UTC()
	var description *string
	if reqContest.Description != "" {
		description = &reqContest.Description
	}

	newContest := models.Contest{
		Title:           reqContest.Title,
		UserPrefix:      reqContest.UserPrefix,
		Description:     description,
		StartTime:       startTime,
		EndTime:         endTime,
		DurationSeconds: reqContest.DurationSeconds,
	}

	err = h.db.Create(&newContest).Error
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create contest", err)
		return
	}

	utils.SendResponse(w, http.StatusCreated, "Contest created successfully", nil)
}
