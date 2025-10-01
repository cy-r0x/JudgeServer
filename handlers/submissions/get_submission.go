package submissions

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
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
	log.Println(userId, submissionId)

	var submission Submission
	err := h.db.Get(&submission, `SELECT * FROM submissions WHERE id=$1`, submissionId)
	if err != nil {
		log.Println("DB Query Error:", err)
		utils.SendResponse(w, http.StatusNotFound, "Submission not found")
		return
	}

	if submission.UserId != userId {
		utils.SendResponse(w, http.StatusForbidden, "Not authorized to view this submission")
		return
	}

	utils.SendResponse(w, http.StatusOK, submission)
}
