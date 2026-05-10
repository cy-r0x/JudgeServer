package submissions

import (
	"log/slog"
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListAllSubmissions(w http.ResponseWriter, r *http.Request) {
	contestId := r.PathValue("contestId")
	if contestId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID required", nil)
		return
	}

	// --- Build params ---
	params := parseSubmissionListParams(r)
	params.ContestID = contestId

	// --- Fetch & Respond ---
	submissions, totalCount, err := h.fetchSubmissions(params)
	if err != nil {
		slog.Error("DB error", "error", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch submissions", nil)
		return
	}

	sendPaginatedResponse(w, submissions, totalCount, params.Limit, params.Page)
}
