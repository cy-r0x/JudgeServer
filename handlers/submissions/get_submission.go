package submissions

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetSubmission(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token")
		return
	}
	userId := payload.Sub
	submissionId := r.PathValue("submissonId")

	var submission models.Submission
	if err := h.db.
		Select("id", "user_id", "username", "problem_id", "contest_id", "language", "source_code", "verdict", "execution_time", "memory_used", "submitted_at").
		Where("id = ?", submissionId).
		First(&submission).Error; err != nil {
		log.Println("DB Query Error:", err)
		utils.SendResponse(w, http.StatusNotFound, "Submission not found")
		return
	}

	if payload.Role != "admin" {
		if submission.UserID != uint(userId) {
			utils.SendResponse(w, http.StatusForbidden, "Not authorized to view this submission")
			return
		}
	}

	utils.SendResponse(w, http.StatusOK, submission)
}
