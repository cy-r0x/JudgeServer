package contest

import (
	"net/http"
	"time"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetContest(w http.ResponseWriter, r *http.Request) {
	contestId := r.PathValue("contestId")
	if contestId == "" {
		utils.SendResponse(w, http.StatusNotFound, "Contest Not Found")
		return
	}

	// Get Contest Information
	var contest Contest
	query := `SELECT id, title, description, start_time, duration_seconds, created_at 
			 FROM contests WHERE id = $1`

	err := h.db.QueryRow(query, contestId).Scan(
		&contest.Id,
		&contest.Title,
		&contest.Description,
		&contest.StartTime,
		&contest.DurationSeconds,
		&contest.CreatedAt,
	)

	if err != nil {
		utils.SendResponse(w, http.StatusNotFound, "Contest Not Found")
		return
	}

	// Calculate contest status
	now := time.Now()
	endTime := contest.StartTime.Add(time.Duration(contest.DurationSeconds) * time.Second)

	if now.Before(contest.StartTime) {
		contest.Status = "UPCOMING"
	} else if now.After(endTime) {
		contest.Status = "ENDED"
	} else {
		contest.Status = "RUNNING"
	}

	// Get Contest Problems
	type Problem struct {
		Id    int64  `json:"id"`
		Title string `json:"title"`
		Slug  string `json:"slug"`
		Index int    `json:"index"`
	}

	problems := []Problem{}

	if contest.Status != "UPCOMING" {
		// First get the contest problems data from contest_problems table
		// then join with problems table to get title and slug
		problemsQuery := `
		SELECT cp.problem_id, p.title, p.slug, cp.index
		FROM contest_problems cp
		JOIN problems p ON cp.problem_id = p.id
		WHERE cp.contest_id = $1
		ORDER BY cp.index
	`

		rows, err := h.db.Query(problemsQuery, contestId)
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contest problems")
			return
		}
		defer rows.Close()

		for rows.Next() {
			var problem Problem
			if err := rows.Scan(&problem.Id, &problem.Title, &problem.Slug, &problem.Index); err != nil {
				utils.SendResponse(w, http.StatusInternalServerError, "Error parsing problem data")
				return
			}
			problems = append(problems, problem)
		}
	}
	// Prepare response with both contest and problems information
	response := struct {
		Contest  Contest   `json:"contest"`
		Problems []Problem `json:"problems"`
	}{
		Contest:  contest,
		Problems: problems,
	}

	utils.SendResponse(w, http.StatusOK, response)
}
