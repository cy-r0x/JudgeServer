package contest

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetContest(w http.ResponseWriter, r *http.Request) {
	var userId *int64
	header := r.Header.Get("Authorization")

	if header != "" {
		headerArr := strings.Split(header, " ")
		if len(headerArr) == 2 {
			accessToken := headerArr[1]
			payload, err := middlewares.DecodeToken(accessToken, h.config.SecretKey)
			if err == nil {
				userId = &payload.Sub
			}
		}
	}

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
		Id           int64  `json:"id"`
		Title        string `json:"title"`
		Slug         string `json:"slug"`
		Index        int    `json:"index"`
		Solved       bool   `json:"solved"`
		Attempted    bool   `json:"attempted"`
		TotalSolvers int    `json:"total_solvers"`
	}

	problems := []Problem{}

	if contest.Status != "UPCOMING" {
		// Optimized query using pre-aggregated stats
		problemsQuery := `
			SELECT 
				cp.problem_id, 
				p.title, 
				p.slug, 
				cp.index,
				COALESCE(cup.is_solved, false) as solved,
				COALESCE(cup.attempt_count > 0, false) as attempted,
				COALESCE(cps.solved_count, 0) as total_solvers
			FROM contest_problems cp
			JOIN problems p ON cp.problem_id = p.id
			LEFT JOIN contest_problem_stats cps ON cps.contest_id = cp.contest_id AND cps.problem_id = cp.problem_id
			LEFT JOIN contest_user_problems cup ON cup.contest_id = cp.contest_id 
				AND cup.problem_id = cp.problem_id 
				AND cup.user_id = $2
			WHERE cp.contest_id = $1
			ORDER BY cp.index
		`

		var rows *sql.Rows
		var err error

		if userId != nil {
			rows, err = h.db.Query(problemsQuery, contestId, *userId)
		} else {
			// For unauthenticated users, pass NULL for user_id
			rows, err = h.db.Query(problemsQuery, contestId, nil)
		}

		if err != nil {
			log.Println("Error fetching contest problems:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contest problems")
			return
		}
		defer rows.Close()

		for rows.Next() {
			var problem Problem
			if err := rows.Scan(&problem.Id, &problem.Title, &problem.Slug, &problem.Index, &problem.Solved, &problem.Attempted, &problem.TotalSolvers); err != nil {
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
