package problem

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) fetchTestcases(problemId string, isSample bool) ([]Testcase, error) {
	query := `
		SELECT id, problem_id, input, expected_output, is_sample, created_at
		FROM testcases
		WHERE problem_id = $1`
	args := []any{problemId}

	if isSample {
		query += ` AND is_sample = TRUE`
	}

	query += ` ORDER BY is_sample DESC, id ASC`

	var testcases []Testcase
	if err := h.db.Select(&testcases, query, args...); err != nil {
		log.Println(err)
		return nil, err
	}

	return testcases, nil
}

func (h *Handler) GetProblem(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
		return
	}

	problemId := r.PathValue("problemId")
	if problemId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Problem ID is required")
		return
	}

	var problem Problem
	isSampleOnly := false

	switch payload.Role {
	case "user":
		// Check if the user has access to this problem through their allowed contest
		var exists bool
		err := h.db.Get(&exists, `
			SELECT EXISTS (
				SELECT 1 FROM contest_problems cp
				WHERE cp.problem_id = $1 AND cp.contest_id = $2
			)
		`, problemId, payload.AllowedContest)

		if err != nil {
			log.Println("Error checking problem access:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to verify problem access")
			return
		}

		if !exists {
			utils.SendResponse(w, http.StatusForbidden, "You don't have access to this problem")
			return
		}

		isSampleOnly = true
		// Fall through to fetch problem data
	case "setter":
		// Check if the problem was created by this setter
		var createdBy int64
		err := h.db.Get(&createdBy, `
			SELECT created_by FROM problems WHERE id = $1
		`, problemId)

		if err != nil {
			log.Println("Error checking problem creator:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to verify problem creator")
			return
		}

		if createdBy != payload.Sub {
			utils.SendResponse(w, http.StatusForbidden, "You don't have access to this problem")
			return
		}

		// Fall through to fetch problem data
	case "admin":
		// Admin has access to all problems
	default:
		utils.SendResponse(w, http.StatusForbidden, "Invalid role")
		return
	}

	// Fetch problem details
	err := h.db.Get(&problem, `
		SELECT id, title, slug, statement, input_statement as input_statement, 
		output_statement, time_limit, memory_limit, created_by, created_at
		FROM problems WHERE id = $1
	`, problemId)

	if err != nil {
		log.Println("Error fetching problem:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch problem")
		return
	}

	// Fetch testcases
	testcases, tcErr := h.fetchTestcases(problemId, isSampleOnly)
	if tcErr != nil {
		log.Println("Error fetching testcases:", tcErr)
		// Continue anyway as we at least have the problem data
	}

	problem.Testcases = testcases

	utils.SendResponse(w, http.StatusOK, problem)
}
