package contest

import (
	"encoding/json"
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) CreateContest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var contest Contest
	err := decoder.Decode(&contest)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Insert Contest into DB
	query := `INSERT INTO contests (title, description, start_time, duration_seconds) 
			 VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	err = h.db.QueryRow(query,
		contest.Title,
		contest.Description,
		contest.StartTime,
		contest.DurationSeconds).Scan(&contest.Id, &contest.CreatedAt)

	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create contest")
		return
	}

	utils.SendResponse(w, http.StatusCreated, contest)
}
