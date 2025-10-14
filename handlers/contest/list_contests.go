package contest

import (
	"net/http"
	"time"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListContests(w http.ResponseWriter, r *http.Request) {
	contests := []Contest{}

	query := `SELECT id, title, start_time, duration_seconds FROM contests ORDER BY start_time DESC`

	rows, err := h.db.Query(query)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contests")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var contest Contest
		if err := rows.Scan(
			&contest.Id,
			&contest.Title,
			&contest.StartTime,
			&contest.DurationSeconds); err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, "Error parsing contest data")
			return
		}

		// Calculate contest status
		now := time.Now()
		endTime := contest.StartTime.Add(time.Duration(contest.DurationSeconds) * time.Second)

		if now.Before(contest.StartTime) {
			contest.Status = "UPCOMING"
		} else if now.After(endTime) {
			contest.Status = "ENDED"
		} else {
			contest.Status = "RUNNING"
		}

		contests = append(contests, contest)
	}

	utils.SendResponse(w, http.StatusOK, contests)
}
