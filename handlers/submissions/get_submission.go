package submissions

import (
	"log"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetSubmission(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token", nil)
		return
	}

	userId := payload.Sub

	submissionIdStr := r.PathValue("submissionId")
	submissionId, err := strconv.ParseInt(submissionIdStr, 10, 64)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid submission ID", nil)
		return
	}

	var submission SubmissionResponse

	err = h.db.
		Table("submissions s").
		Select(`
			s.id,
			s.user_id,
			s.problem_id,
			s.contest_id,
			s.language,
			s.source_code,
			s.status,
			s.exec_time,
			s.exec_memory,
			s.created_at,
			u.name as user_name,
			u.username as username
		`).
		Joins("JOIN users u ON u.id = s.user_id").
		Where("s.id = ?", submissionId).
		Scan(&submission).Error

	if err != nil {
		log.Println("DB Query Error:", err)
		utils.SendResponse(w, http.StatusNotFound, "Submission not found", nil)
		return
	}

	// authorization check
	if payload.Role != "admin" {
		if submission.UserId != userId {
			utils.SendResponse(w, http.StatusForbidden, "Not authorized to view this submission", nil)
			return
		}
	}

	utils.SendResponse(w, http.StatusOK, "Submission Fetched Successfully", submission)
}
