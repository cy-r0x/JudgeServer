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
