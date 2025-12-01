package submissions

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListUserSubmissions(w http.ResponseWriter, r *http.Request) {

	limit := 20
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	strLimit := r.URL.Query().Get("limit")
	verdictFilter := r.URL.Query().Get("verdict")
	fmt.Println(verdictFilter)

	if strLimit != "" {
		parsedLimit, err := strconv.Atoi(strLimit)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	crrPage, err := strconv.Atoi(page)
	if err != nil {
		log.Println(err)
		utils.SendResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	userId := payload.Sub
	contestId := payload.AllowedContest

	if contestId == nil {
		utils.SendResponse(w, http.StatusBadRequest, "No contest specified")
		return
	}

	offset := (crrPage - 1) * limit

	// Optimized single query with window function for count
	type SubmissionWithCount struct {
		Submission
		TotalCount int `db:"total_count"`
	}

	var results []SubmissionWithCount
	query := `
		SELECT 
			id, username, problem_id, language, 
			verdict, execution_time, memory_used, submitted_at,
			COUNT(*) OVER() as total_count
		FROM submissions 
		WHERE user_id=$1 AND contest_id=$2 
	`

	if verdictFilter != "" {
		query += " AND verdict=$3"
		query += " ORDER BY submitted_at DESC LIMIT $4 OFFSET $5"
		err = h.db.Select(&results, query, userId, *contestId, verdictFilter, limit, offset)
	} else {
		query += " ORDER BY submitted_at DESC LIMIT $3 OFFSET $4"
		err = h.db.Select(&results, query, userId, *contestId, limit, offset)
	}

	if err != nil {
		log.Println("DB Query Error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch submissions")
		return
	}

	var submissions []Submission
	totalCount := 0
	for _, r := range results {
		submissions = append(submissions, r.Submission)
		totalCount = r.TotalCount
	}

	response := struct {
		Submissions []Submission `json:"submissions"`
		TotalItem   int          `json:"total_item"`
		TotalPages  int          `json:"total_pages"`
		Limit       int          `json:"limit"`
		Page        int          `json:"page"`
	}{
		Submissions: submissions,
		TotalItem:   totalCount,
		TotalPages:  (totalCount + limit - 1) / limit,
		Limit:       limit,
		Page:        crrPage,
	}

	utils.SendResponse(w, http.StatusOK, response)
}
