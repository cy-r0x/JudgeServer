package submissions

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListAllSubmissions(w http.ResponseWriter, r *http.Request) {
	contestId := r.PathValue("contestId")
	log.Println(contestId)

	var submissions []Submission
	err := h.db.Select(&submissions, `
		SELECT id, user_id, username, problem_id, contest_id, language, 
		       verdict, execution_time, memory_used, submitted_at, first_blood
		FROM submissions 
		WHERE contest_id=$1 
		ORDER BY submitted_at DESC
	`, contestId)
	if err != nil {
		log.Println("DB Query Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch submissions")
		return
	}

	utils.SendResponse(w, http.StatusOK, submissions)
}
