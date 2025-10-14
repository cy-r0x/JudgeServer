package contest

import (
	"encoding/json"
	"net/http"

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

	query := `UPDATE contests 
			 SET title = $1, description = $2, start_time = $3, duration_seconds = $4
			 WHERE id = $5 RETURNING created_at`

	err = h.db.QueryRow(query,
		contest.Title,
		contest.Description,
		contest.StartTime,
		contest.DurationSeconds,
		contest.Id).Scan(&contest.CreatedAt)

	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update contest")
		return
	}

	utils.SendResponse(w, http.StatusOK, contest)
}
