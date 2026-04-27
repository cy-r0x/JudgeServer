package contest

import (
	"net/http"
	"time"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListContests(w http.ResponseWriter, r *http.Request) {
	var contests []Contest
	err := h.db.Order("start_time DESC").Find(&contests).Error
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contests", nil)
		return
	}

	now := time.Now()
	for i := range contests {
		if now.Before(contests[i].StartTime) {
			contests[i].Status = "UPCOMING"
		} else if now.After(contests[i].EndTime) {
			contests[i].Status = "ENDED"
		} else {
			contests[i].Status = "RUNNING"
		}
	}

	if contests == nil {
		contests = []Contest{}
	}

	utils.SendResponse(w, http.StatusOK, "Contests fetched successfully", contests)
}
