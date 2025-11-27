package contest_problems

import (
	"log"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetContestProblems(w http.ResponseWriter, r *http.Request) {
	contestId := r.PathValue("contestId")

	// Convert contestId to int64
	contestIdInt, err := strconv.ParseInt(contestId, 10, 64)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid contest ID")
		return
	}

	// Check if contest exists
	var contestExists bool
	if err = h.db.Get(&contestExists, `SELECT EXISTS(SELECT 1 FROM contests WHERE id=$1)`, contestIdInt); err != nil {
		log.Println("Failed to check contest existence:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to get contest problems")
		return
	}
	if !contestExists {
		utils.SendResponse(w, http.StatusNotFound, "Contest does not exist")
		return
	}

	// Get all contest problems with problem details and author information in one optimized query
	var contestProblems []ContestProblem
	query := `
		SELECT 
			cp.contest_id,
			cp.problem_id,
			cp.index,
			p.title as problem_name,
			COALESCE(u.full_name, 'Unknown') as problem_author
		FROM contest_problems cp
		JOIN problems p ON cp.problem_id = p.id
		LEFT JOIN users u ON p.created_by = u.id
		WHERE cp.contest_id = $1
		ORDER BY cp.index ASC
	`

	if err = h.db.Select(&contestProblems, query, contestIdInt); err != nil {
		log.Println("Failed to get contest problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to get contest problems")
		return
	}

	utils.SendResponse(w, http.StatusOK, contestProblems)
}
