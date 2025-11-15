package submissions

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListAllSubmissions(w http.ResponseWriter, r *http.Request) {
	const limit = 10

	contestId := r.PathValue("contestId")
	log.Println(contestId)

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
		query := `SELECT COUNT(*) FROM submissions WHERE contest_id=$1`
		err := h.db.Get(&response.TotalItem, query, contestId)
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
		query := `SELECT sub.id, sub.user_id, u.username, sub.problem_id, sub.contest_id, sub.language, 
		       sub.verdict, sub.execution_time, sub.memory_used, sub.submitted_at, u.clan, u.full_name, u.room_no, u.pc_no
		FROM submissions sub LEFT JOIN users u ON sub.user_id = u.id
		WHERE contest_id=$1 
		ORDER BY submitted_at DESC LIMIT $2 OFFSET $3`
		err := h.db.Select(&submissions, query, contestId, limit, offset)
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
