package submissions

import (
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
	if submission.ProblemId == "" || submission.ContestId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Problem ID and Contest ID are required")
		return
	}

	// Check if problem exists
	var problemExists bool
	if err := h.db.Model(&models.Problem{}).Select("1").Where("id = ?", submission.ProblemId).Find(&problemExists).Error; err != nil {
		log.Println("Failed to check problem existence:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to validate submission")
		return
	}
	if !problemExists {
		utils.SendResponse(w, http.StatusBadRequest, "Problem does not exist")
		return
	}

	// Check if contest exists and get contest timing information
	var contest models.Contest
	if err := h.db.Select("start_time", "duration_seconds").Where("id = ?", submission.ContestId).First(&contest).Error; err != nil {
		log.Println("Failed to get contest details:", err)
		utils.SendResponse(w, http.StatusBadRequest, "Contest does not exist")
		return
	}

	// Check if contest is currently running
	now := time.Now()
	endTime := contest.StartTime.Add(time.Duration(contest.DurationSeconds) * time.Second)

	if now.Before(contest.StartTime) {
		utils.SendResponse(w, http.StatusBadRequest, "Contest has not started yet")
		return
	}

	if now.After(endTime) {
		utils.SendResponse(w, http.StatusBadRequest, "Contest has ended")
		return
	}

	var problem Problem
	if err := h.db.Model(&models.Problem{}).
		Select("time_limit", "memory_limit", "checker_type", "checker_strict_space", "checker_precision").
		Where("id = ?", submission.ProblemId).
		First(&problem).Error; err != nil {
		log.Println("Failed to fetch problem details:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to validate submission")
		return
	}

	var testcases []Testcase
	if err := h.db.Model(&models.Testcase{}).
		Select("input", "expected_output").
		Where("problem_id = ?", submission.ProblemId).
		Find(&testcases).Error; err != nil {
		log.Println("Failed to fetch testcases:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to validate submission")
		return
	}

	problem.Testcases = testcases

	// Check if problem is assigned to the contest
	var problemInContest bool
	if err := h.db.Model(&models.ContestProblem{}).Select("1").Where("contest_id = ? AND problem_id = ?", submission.ContestId, submission.ProblemId).Find(&problemInContest).Error; err != nil {
		log.Println("Failed to check problem assignment:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to validate submission")
		return
	}
	if !problemInContest {
		utils.SendResponse(w, http.StatusBadRequest, "Problem is not assigned to this contest")
		return
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	contestID := submission.ContestId
	newSubmission := models.Submission{
		UserID:     userId,
		Username:   username,
		ProblemID:  submission.ProblemId,
		ContestID:  &contestID,
		Language:   submission.Language,
		SourceCode: submission.SourceCode,
	}

	if err := tx.Create(&newSubmission).Error; err != nil {
		tx.Rollback()
		log.Println("DB Insert Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create submission")
		return
	}

	problem.SubmissionId = newSubmission.ID
	problem.SourceCode = submission.SourceCode
	problem.Language = submission.Language

	if err := h.submitToQueue(&problem); err != nil {
		tx.Rollback()
		log.Println("Queue Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to enqueue submission")
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Commit Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	utils.SendResponse(w, http.StatusOK, map[string]any{"submission_id": problem.SubmissionId})
}
