package submissions

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	const maxBodySize = 50 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token", nil)
		return
	}
	userID := payload.Sub

	var submission UserSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if submission.ProblemId == "" || submission.ContestId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Problem ID and Contest ID are required", nil)
		return
	}

	now := time.Now()

	// -----------------------------------
	// 1. Contest validation
	// -----------------------------------
	var startTime, endTime time.Time

	err := h.db.Raw(`
		SELECT start_time, end_time
		FROM contests
		WHERE id = $1
	`, submission.ContestId).Row().Scan(&startTime, &endTime)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendResponse(w, http.StatusBadRequest, "Contest does not exist", nil)
			return
		}
		log.Println("Contest query error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to validate submission", nil)
		return
	}

	if now.Before(startTime) {
		utils.SendResponse(w, http.StatusBadRequest, "Contest has not started yet", nil)
		return
	}
	if now.After(endTime) {
		utils.SendResponse(w, http.StatusBadRequest, "Contest has ended", nil)
		return
	}

	// -----------------------------------
	// 2. Problem + contest assignment
	// -----------------------------------
	var problem models.Problem

	err = h.db.Raw(`
		SELECT 
			p.id,
			p.time_limit,
			p.memory_limit,
			p.checker_type,
			p.checker_strict_space,
			p.checker_precision
		FROM problems p
		JOIN contest_problems cp 
			ON cp.problem_id = p.id
		WHERE 
			p.id = $1
			AND cp.contest_id = $2
		LIMIT 1
	`, submission.ProblemId, submission.ContestId).
		Scan(&problem).Error

	if err != nil {
		log.Println("Problem query error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to validate submission", nil)
		return
	}

	if problem.Id == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Problem not found in this contest", nil)
		return
	}

	// -----------------------------------
	// 3. Fetch testcases
	// -----------------------------------
	var testcases []Testcase

	err = h.db.Raw(`
		SELECT input, expected_output
		FROM testcases
		WHERE problem_id = $1
	`, submission.ProblemId).
		Scan(&testcases).Error

	if err != nil {
		log.Println("Testcase query error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to validate submission", nil)
		return
	}

	var queueSubmission QueueSubmission

	queueSubmission.Testcases = testcases

	// -----------------------------------
	// 4. Insert submission (transaction)
	// -----------------------------------
	tx := h.db.Begin()
	if tx.Error != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to start transaction", nil)
		return
	}

	var submissionID int64

	err = tx.Raw(`
		INSERT INTO submissions (
			user_id,
			problem_id,
			contest_id,
			language,
			source_code
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		userID,
		submission.ProblemId,
		submission.ContestId,
		submission.Language,
		submission.SourceCode,
	).Row().Scan(&submissionID)

	if err != nil {
		tx.Rollback()
		log.Println("Insert error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create submission", nil)
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Commit error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to commit transaction", nil)
		return
	}

	// -----------------------------------
	// 5. Send to judge queue
	// -----------------------------------
	queueSubmission.SubmissionId = submissionID
	queueSubmission.SourceCode = submission.SourceCode
	queueSubmission.Language = submission.Language

	if err := h.submitToQueue(&queueSubmission); err != nil {
		// IMPORTANT: do NOT rollback (already committed)
		log.Println("Queue error:", err)

		// Optional: mark submission as failed in DB
		utils.SendResponse(w, http.StatusInternalServerError, "Submission stored but failed to queue", nil)
		return
	}

	// -----------------------------------
	// DONE
	// -----------------------------------
	utils.SendResponse(w, http.StatusOK, "Submission Created Successfully", nil)
}
