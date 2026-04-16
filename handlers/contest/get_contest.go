package contest

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetContest(w http.ResponseWriter, r *http.Request) {
	var userId *string
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
	var dbContest models.Contest
	err := h.db.
		Select("id", "title", "description", "start_time", "duration_seconds", "created_at").
		Where("id = ?", contestId).
		First(&dbContest).Error

	if err != nil {
		utils.SendResponse(w, http.StatusNotFound, "Contest Not Found")
		return
	}

	description := ""
	if dbContest.Description != nil {
		description = *dbContest.Description
	}

	contest := Contest{
		Id:              dbContest.ID,
		Title:           dbContest.Title,
		Description:     description,
		StartTime:       dbContest.StartTime,
		DurationSeconds: dbContest.DurationSeconds,
		CreatedAt:       dbContest.CreatedAt,
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
		Id           string `json:"id"`
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
		// Note: GORM Raw querying is used for complex joins with COALESCE and dynamic parameter checking.
		problemsQuery := `
			SELECT 
				cp.problem_id as id, 
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
				AND cup.user_id = @userId
			WHERE cp.contest_id = @contestId
			ORDER BY cp.index
		`

		var uId sql.NullString
		if userId != nil {
			uId = sql.NullString{String: *userId, Valid: true}
		} else {
			uId = sql.NullString{Valid: false}
		}

		err = h.db.Raw(problemsQuery, sql.Named("contestId", contestId), sql.Named("userId", uId)).Scan(&problems).Error
		if err != nil {
			log.Println("Error fetching contest problems:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contest problems")
			return
		}
	}

	if problems == nil {
		problems = []Problem{}
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
