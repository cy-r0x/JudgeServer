package submissions

import (
	"log"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListAllSubmissions(w http.ResponseWriter, r *http.Request) {
	var limit int = 20

	contestId := r.PathValue("contestId")
	strLimit := r.URL.Query().Get("limit")
	verdictFilter := r.URL.Query().Get("verdict")

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

	offset := (crrPage - 1) * limit

	// Optimized single query with window function for count
	type SubmissionWithCount struct {
		Submission
		TotalCount int `db:"total_count"`
	}

	var results []SubmissionWithCount
	query := `
		SELECT 
			sub.id, sub.user_id, u.username, sub.problem_id, sub.contest_id, sub.language, 
			sub.verdict, sub.execution_time, sub.memory_used, sub.submitted_at, 
			u.clan, u.full_name, u.room_no, u.pc_no,
			COUNT(*) OVER() as total_count
		FROM submissions sub 
		LEFT JOIN users u ON sub.user_id = u.id
		WHERE sub.contest_id=$1
	`

	if verdictFilter != "" {
		query += " AND sub.verdict=$2"
		query += " ORDER BY sub.submitted_at DESC LIMIT $3 OFFSET $4"
	} else {
		query += " ORDER BY sub.submitted_at DESC LIMIT $2 OFFSET $3"
	}

	if verdictFilter != "" {
		err = h.db.Select(&results, query, contestId, verdictFilter, limit, offset)
	} else {
		err = h.db.Select(&results, query, contestId, limit, offset)
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
