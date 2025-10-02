package submissions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token")
		return
	}
	userId := payload.Sub
	username := payload.Username

	var submission UserSubmission
	if err := decoder.Decode(&submission); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if submission.ProblemId == 0 || submission.ContestId == 0 {
		utils.SendResponse(w, http.StatusBadRequest, "Problem ID and Contest ID are required")
		return
	}

	// Check if problem exists
	var problemExists bool
	if err := h.db.Get(&problemExists, `SELECT EXISTS(SELECT 1 FROM problems WHERE id=$1)`, submission.ProblemId); err != nil {
		log.Println("Failed to check problem existence:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to validate submission")
		return
	}
	if !problemExists {
		utils.SendResponse(w, http.StatusBadRequest, "Problem does not exist")
		return
	}

	// Check if contest exists
	var contestExists bool
	if err := h.db.Get(&contestExists, `SELECT EXISTS(SELECT 1 FROM contests WHERE id=$1)`, submission.ContestId); err != nil {
		log.Println("Failed to check contest existence:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to validate submission")
		return
	}
	if !contestExists {
		utils.SendResponse(w, http.StatusBadRequest, "Contest does not exist")
		return
	}

	// Check if problem is assigned to the contest
	var problemInContest bool
	if err := h.db.Get(&problemInContest, `SELECT EXISTS(SELECT 1 FROM contest_problems WHERE contest_id=$1 AND problem_id=$2)`, submission.ContestId, submission.ProblemId); err != nil {
		log.Println("Failed to check problem assignment:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to validate submission")
		return
	}
	if !problemInContest {
		utils.SendResponse(w, http.StatusBadRequest, "Problem is not assigned to this contest")
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	query := `INSERT INTO submissions (user_id, username, problem_id, contest_id, language, source_code) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var submissionId int64
	err = tx.QueryRow(query, userId, username, submission.ProblemId, submission.ContestId, submission.Language, submission.SourceCode).Scan(&submissionId)
	if err != nil {
		tx.Rollback()
		log.Println("DB Insert Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create submission")
		return
	}

	err = h.submitToQueue(submissionId, &submission)
	if err != nil {
		tx.Rollback()
		log.Println("Queue Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to enqueue submission")
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Commit Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	utils.SendResponse(w, http.StatusOK, map[string]any{"submission_id": submissionId})
}
