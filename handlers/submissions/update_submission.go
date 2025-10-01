package submissions

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) UpdateSubmission(w http.ResponseWriter, r *http.Request) {
	engineData, ok := r.Context().Value("engineData").(*middlewares.EngineData)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token")
		return
	}
	log.Println(engineData)

	// Update the submission in the DB
	query := `UPDATE submissions SET verdict=$1, execution_time=$2, memory_used=$3 WHERE id=$4`
	_, err := h.db.Exec(query, engineData.Verdict, engineData.ExecutionTime, engineData.ExecutionMemory, engineData.SubmissionId)
	if err != nil {
		log.Println("DB Update Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update submission")
		return
	}

	if isAcceptedVerdict(engineData.Verdict) {
		h.updateStandingsForAccepted(engineData.SubmissionId)
	}

	utils.SendResponse(w, http.StatusOK, map[string]any{"message": "Submission updated"})
}
