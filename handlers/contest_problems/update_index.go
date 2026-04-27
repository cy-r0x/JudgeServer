package contest_problems

import (
	"encoding/json"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) UpdateContestIndex(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var contestProblems []ContestProblem
	if err := decoder.Decode(&contestProblems); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to start transaction", nil)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	totalRowsAffected := int64(0)

	for _, contestProblem := range contestProblems {
		result := tx.Model(&models.ContestProblem{}).
			Where("contest_id = ? AND problem_id = ?", contestProblem.ContestId, contestProblem.ProblemId).
			Update("index", contestProblem.Index)

		if result.Error != nil {
			tx.Rollback()
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to update problem index", nil)
			return
		}

		totalRowsAffected += result.RowsAffected
	}

	if totalRowsAffected == 0 {
		tx.Rollback()
		utils.SendResponse(w, http.StatusNotFound, "No contest problems updated", nil)
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to commit transaction", nil)
		return
	}

	utils.SendResponse(w, http.StatusOK, contestProblems, nil)
}
