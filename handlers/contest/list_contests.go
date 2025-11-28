package contest

import (
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListContests(w http.ResponseWriter, r *http.Request) {
	var contests []Contest

	// Use CASE for status calculation in SQL for better performance
	query := `
		SELECT 
			id, title, start_time, duration_seconds,
			CASE 
				WHEN start_time > NOW() THEN 'UPCOMING'
				WHEN start_time + (duration_seconds || ' seconds')::INTERVAL < NOW() THEN 'ENDED'
				ELSE 'RUNNING'
			END as status
		FROM contests 
		ORDER BY start_time DESC
	`

	type ContestWithStatus struct {
		Contest
		Status string `db:"status"`
	}

	var results []ContestWithStatus
	err := h.db.Select(&results, query)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contests")
		return
	}

	for _, r := range results {
		contest := r.Contest
		contest.Status = r.Status
		contests = append(contests, contest)
	}

	utils.SendResponse(w, http.StatusOK, contests)
}
