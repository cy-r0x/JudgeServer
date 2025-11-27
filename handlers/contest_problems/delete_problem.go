package contest_problems

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) DeleteContestProblem(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var contestProblem ContestProblem
	err := decoder.Decode(&contestProblem)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if contestProblem.ContestId == 0 || contestProblem.ProblemId == 0 {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID and Problem ID are required")
		return
	}

	tx, err := h.db.Beginx()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem")
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, execErr := tx.Exec(`DELETE FROM contest_problems WHERE contest_id=$1 AND problem_id=$2`, contestProblem.ContestId, contestProblem.ProblemId)
	if execErr != nil {
		err = execErr
		log.Println("Failed to delete contest problem:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback()
		utils.SendResponse(w, http.StatusNotFound, "Problem not assigned to this contest")
		return
	}

	var remaining []ContestProblem
	if err = tx.Select(&remaining, `SELECT contest_id, problem_id, index FROM contest_problems WHERE contest_id=$1 ORDER BY index ASC`, contestProblem.ContestId); err != nil {
		log.Println("Failed to fetch remaining contest problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem")
		return
	}

	for i, cp := range remaining {
		newIndex := i + 1
		if cp.Index == newIndex {
			continue
		}
		if _, err = tx.Exec(`UPDATE contest_problems SET index=$1 WHERE contest_id=$2 AND problem_id=$3`, newIndex, cp.ContestId, cp.ProblemId); err != nil {
			log.Println("Failed to update contest problem index:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem")
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Println("Failed to commit contest problem deletion:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem")
		return
	}

	utils.SendResponse(w, http.StatusOK, map[string]any{"message": "Contest problem removed"})
}
