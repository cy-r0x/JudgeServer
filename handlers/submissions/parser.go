package submissions

import (
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/utils"
)

// Parse query params (same for both endpoints)
func parseSubmissionListParams(r *http.Request) SubmissionListParams {
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	return SubmissionListParams{
		Limit:          limit,
		Page:           page,
		Status:         r.URL.Query().Get("status"),
		SearchName:     r.URL.Query().Get("searchName"),
		SearchUsername: r.URL.Query().Get("searchUsername"),
		// ContestID and UserID set by caller
	}
}

// Standardized JSON response
func sendPaginatedResponse(w http.ResponseWriter, submissions []SubmissionResponse, totalCount, limit, page int) {
	totalPages := 0
	if totalCount > 0 {
		totalPages = (totalCount + limit - 1) / limit
	}

	utils.SendResponse(w, http.StatusOK, "Submissions fetched successfully", map[string]any{
		"submissions": submissions,
		"totalItem":   totalCount,
		"totalPages":  totalPages,
		"limit":       limit,
		"page":        page,
	})
}
