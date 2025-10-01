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
	err := h.db.Select(&submissions, `SELECT * FROM submissions WHERE contest_id=$1 ORDER BY submitted_at DESC`, contestId)
	if err != nil {
		log.Println("DB Query Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch submissions")
		return
	}

	utils.SendResponse(w, http.StatusOK, submissions)
}
