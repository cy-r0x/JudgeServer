package submissions

import (
	"log/slog"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) UpdateSubmission(w http.ResponseWriter, r *http.Request) {
	enginePayload, ok := r.Context().Value("enginePayload").(*middlewares.EnginePayload)

	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token", nil)
		return
	}

	// -----------------------------
	// Validate status (important)
	// -----------------------------
	validStatuses := map[string]bool{
		"ACCEPTED":              true,
		"WRONG_ANSWER":          true,
		"TIME_LIMIT_EXCEEDED":   true,
		"RUNTIME_ERROR":         true,
		"MEMORY_LIMIT_EXCEEDED": true,
		"COMPILATION_ERROR":     true,
		"INTERNAL_ERROR":        true,
	}

	if !validStatuses[enginePayload.Status] {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid status from engine", nil)
		return
	}

	// -----------------------------
	// Build update map
	// -----------------------------
	updates := map[string]interface{}{
		"status": enginePayload.Status,
	}

	if enginePayload.ExecutionTime != nil {
		updates["exec_time"] = *enginePayload.ExecutionTime
	}

	if enginePayload.ExecutionMemory != nil {
		updates["exec_memory"] = *enginePayload.ExecutionMemory
	}

	// -----------------------------
	// Update submission
	// -----------------------------
	result := h.db.Model(&models.Submission{}).
		Where("id = ?", enginePayload.SubmissionId).
		Updates(updates)

	if result.Error != nil {
		slog.Error("DB Update Error", "submission_id", enginePayload.SubmissionId, "error", result.Error)

		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update submission", nil)
		return
	}

	if result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusNotFound, "Submission not found", nil)
		return
	}

	// -----------------------------
	// Standings update (sync call)
	// -----------------------------
	if enginePayload.Status == "ACCEPTED" {
		if err := h.updateStandingsForAccepted(enginePayload.SubmissionId); err != nil {
			slog.Error("Standings update (ACCEPTED) failed", "error", err)
		}
	} else {
		if err := h.updateStandingsForNonAccepted(
			enginePayload.SubmissionId,
			enginePayload.Status,
		); err != nil {
			slog.Error("Standings update failed", "error", err)
		}
	}

	// -----------------------------
	// Done
	// -----------------------------
	utils.SendResponse(w, http.StatusOK, "Submission updated", nil)
}
