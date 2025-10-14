package submissions

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListUserSubmissions(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token")
		return
	}
	userId := payload.Sub
	contestId := payload.AllowedContest
	log.Println(userId, contestId)

	if contestId == nil {
		utils.SendResponse(w, http.StatusBadRequest, "No contest specified")
		return
	}

	var submissions []Submission
	err := h.db.Select(&submissions, `
		SELECT id, user_id, username, problem_id, contest_id, language, 
		       verdict, execution_time, memory_used, submitted_at, first_blood
		FROM submissions 
		WHERE user_id=$1 AND contest_id=$2 
		ORDER BY submitted_at DESC
	`, userId, *contestId)
	if err != nil {
		log.Println("DB Query Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch submissions")
		return
	}

	utils.SendResponse(w, http.StatusOK, submissions)
}
