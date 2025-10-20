package submissions

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListUserSubmissions(w http.ResponseWriter, r *http.Request) {

	const limit = 10

	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token")
		return
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

	var submissions []Submission
	var response struct {
		Submissions []Submission `json:"submissions"`
		TotalItem   int          `json:"total_item"`
		TotalPages  int          `json:"total_pages"`
		Limit       int          `json:"limit"`
		Page        int          `json:"page"`
	}

	response.Limit = limit
	response.Page = crrPage
	response.TotalItem = 0

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		query := `SELECT COUNT(*) FROM submissions WHERE user_id=$1 AND contest_id=$2`
		err := h.db.Get(&response.TotalItem, query, userId, *contestId)
		if err != nil {
			log.Println("DB Count Error:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Database Error")
			return
		}
		response.TotalPages = (response.TotalItem + limit - 1) / limit
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		offset := (crrPage - 1) * limit
		query := `SELECT id, username, problem_id, language, 
		       verdict, execution_time, memory_used, submitted_at
		FROM submissions 
		WHERE user_id=$1 AND contest_id=$2 
		ORDER BY submitted_at DESC LIMIT $3 OFFSET $4`
		err := h.db.Select(&submissions, query, userId, *contestId, limit, offset)
		if err != nil {
			log.Println("DB Query Error:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch submissions")
			return
		}
	}()

	wg.Wait()
	response.Submissions = submissions

	utils.SendResponse(w, http.StatusOK, response)
}
