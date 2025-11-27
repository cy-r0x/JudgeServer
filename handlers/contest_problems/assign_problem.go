package contest_problems

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) AssignContestProblems(w http.ResponseWriter, r *http.Request) {
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

	// Check if contest exists
	var contestExists bool
	if err = h.db.Get(&contestExists, `SELECT EXISTS(SELECT 1 FROM contests WHERE id=$1)`, contestProblem.ContestId); err != nil {
		log.Println("Failed to check contest existence:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}
	if !contestExists {
		utils.SendResponse(w, http.StatusBadRequest, "Contest does not exist")
		return
	}

	// Get problem details and author information in a single optimized query
	type ProblemDetails struct {
		Title    string `db:"title"`
		FullName string `db:"full_name"`
	}

	var problemDetails ProblemDetails
	query := `
		SELECT p.title, u.full_name 
		FROM problems p 
		LEFT JOIN users u ON p.created_by = u.id 
		WHERE p.id = $1
	`

	if err = h.db.Get(&problemDetails, query, contestProblem.ProblemId); err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.SendResponse(w, http.StatusBadRequest, "Problem does not exist")
			return
		}
		log.Println("Failed to get problem details:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}

	// Set problem details in the response struct
	contestProblem.ProblemName = problemDetails.Title
	contestProblem.ProblemAuthor = problemDetails.FullName

	tx, err := h.db.Beginx()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Ensure the problem is not already assigned
	var exists bool
	if err = tx.Get(&exists, `SELECT EXISTS(SELECT 1 FROM contest_problems WHERE contest_id=$1 AND problem_id=$2)`, contestProblem.ContestId, contestProblem.ProblemId); err != nil {
		log.Println("Failed to check existing assignment:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}
	if exists {
		tx.Rollback()
		utils.SendResponse(w, http.StatusConflict, "Problem already assigned to this contest")
		return
	}

	// Count existing problems to determine next index
	var count int
	if err = tx.Get(&count, `SELECT COUNT(*) FROM contest_problems WHERE contest_id=$1`, contestProblem.ContestId); err != nil {
		log.Println("Failed to count contest problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}

	contestProblem.Index = count + 1

	if _, err = tx.Exec(`INSERT INTO contest_problems (contest_id, problem_id, index) VALUES ($1, $2, $3)`, contestProblem.ContestId, contestProblem.ProblemId, contestProblem.Index); err != nil {
		log.Println("Failed to insert contest problem:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println("Failed to commit contest problem assignment:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}

	utils.SendResponse(w, http.StatusOK, contestProblem)
}
