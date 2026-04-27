package submissions

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListUserSubmissions(w http.ResponseWriter, r *http.Request) {
	// --- Auth ---
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid token", nil)
		return
	}

	userID := payload.Sub
	contestId := payload.AllowedContest
	if contestId == nil {
		utils.SendResponse(w, http.StatusBadRequest, "No contest specified", nil)
		return
	}

	// --- Build params ---
	params := parseSubmissionListParams(r)
	params.ContestID = *contestId
	params.UserID = &userID

	// --- Fetch & Respond ---
	submissions, totalCount, err := h.fetchSubmissions(params)
	if err != nil {
		log.Println("DB error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch submissions", nil)
		return
	}

	sendPaginatedResponse(w, submissions, totalCount, params.Limit, params.Page)
}
