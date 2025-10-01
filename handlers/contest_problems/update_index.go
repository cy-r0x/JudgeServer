package contest_problems

import (
	"encoding/json"
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) UpdateContestIndex(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var contestProblems []ContestProblem
	if err := decoder.Decode(&contestProblems); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	query := `UPDATE contest_problems 
	          SET index = $1
	          WHERE contest_id = $2 AND problem_id = $3`

	totalRowsAffected := int64(0)

	for _, contestProblem := range contestProblems {
		result, err := tx.Exec(query,
			contestProblem.Index,
			contestProblem.ContestId,
			contestProblem.ProblemId)
		if err != nil {
			tx.Rollback()
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to update problem index")
			return
		}

		rows, _ := result.RowsAffected()
		totalRowsAffected += rows
	}

	if totalRowsAffected == 0 {
		tx.Rollback()
		utils.SendResponse(w, http.StatusNotFound, "No contest problems updated")
		return
	}

	if err := tx.Commit(); err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	utils.SendResponse(w, http.StatusOK, contestProblems)
}
