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
		log.Println("kichu ekta")
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	// Handle nullable execution time and memory values
	var executionTime any
	var memoryUsed any

	if engineData.ExecutionTime != nil {
		executionTime = *engineData.ExecutionTime
	}

	if engineData.ExecutionMemory != nil {
		memoryUsed = *engineData.ExecutionMemory
	}

	// Update the submission in the DB
	query := `UPDATE submissions SET verdict=$1, execution_time=$2, memory_used=$3 WHERE id=$4`
	_, err := h.db.Exec(query, engineData.Verdict, executionTime, memoryUsed, engineData.SubmissionId)
	if err != nil {
		log.Println("DB Update Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update submission")
		return
	}

	if engineData.Verdict == "ac" {
		h.updateStandingsForAccepted(engineData.SubmissionId)
	}

	utils.SendResponse(w, http.StatusOK, map[string]any{"message": "Submission updated"})
}
