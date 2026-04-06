package submissions

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) UpdateSubmission(w http.ResponseWriter, r *http.Request) {
	engineData, ok := r.Context().Value("engineData").(*middlewares.EngineData)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	// Handle nullable execution time and memory values
	var executionTime interface{}
	var memoryUsed interface{}

	if engineData.ExecutionTime != nil {
		executionTime = *engineData.ExecutionTime
	}

	if engineData.ExecutionMemory != nil {
		memoryUsed = *engineData.ExecutionMemory
	}

	updates := map[string]interface{}{
		"verdict": engineData.Verdict,
	}
	if executionTime != nil {
		updates["execution_time"] = executionTime
	}
	if memoryUsed != nil {
		updates["memory_used"] = memoryUsed
	}

	// Update the submission in the DB
	if err := h.db.Model(&models.Submission{}).Where("id = ?", engineData.SubmissionId).Updates(updates).Error; err != nil {
		log.Println("DB Update Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update submission")
		return
	}

	if engineData.Verdict == "ac" {
		h.updateStandingsForAccepted(engineData.SubmissionId)
	} else {
		// Track non-AC submissions for contest standings
		h.updateStandingsForNonAccepted(engineData.SubmissionId, engineData.Verdict)
	}

	utils.SendResponse(w, http.StatusOK, map[string]interface{}{"message": "Submission updated"})
}
