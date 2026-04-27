package contest

import (
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
		utils.SendResponse(w, http.StatusNotFound, "Contest Not Found", nil)
		return
	}

	var dbContest models.Contest
	err := h.db.Where("id = ?", contestId).First(&dbContest).Error
	if err != nil {
		utils.SendResponse(w, http.StatusNotFound, "Contest Not Found", nil)
		return
	}

	description := ""
	if dbContest.Description != nil {
		description = *dbContest.Description
	}

	contest := Contest{
		Id:              dbContest.Id,
		Title:           dbContest.Title,
		UserPrefix:      dbContest.UserPrefix,
		Description:     description,
		StartTime:       dbContest.StartTime,
		EndTime:         dbContest.EndTime,
		DurationSeconds: dbContest.DurationSeconds,
		CreatedAt:       dbContest.CreatedAt,
		UpdatedAt:       dbContest.UpdatedAt,
	}

	// Calculate contest status
	now := time.Now()
	if now.Before(contest.StartTime) {
		contest.Status = "UPCOMING"
	} else if now.After(contest.EndTime) {
		contest.Status = "ENDED"
	} else {
		contest.Status = "RUNNING"
	}

	// Get Contest Problems
	type Problem struct {
		Id           string `json:"id"`
		Title        string `json:"title"`
		Index        int    `json:"index"`
		Solved       bool   `json:"solved"`
		Attempted    bool   `json:"attempted"`
		TotalSolvers int    `json:"total_solvers"`
	}

	problems := []Problem{}

	if contest.Status != "UPCOMING" {
		var contestProblems []models.ContestProblem
		if err := h.db.Where("contest_id = ?", contestId).Order("\"index\" ASC").Find(&contestProblems).Error; err != nil {
			log.Println("Error fetching contest problems:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contest problems", nil)
			return
		}

		var problemResults []models.ContestProblemResult
		if userId != nil {
			if err := h.db.Where("contest_id = ? AND user_id = ?", contestId, *userId).Find(&problemResults).Error; err != nil {
				log.Println("Error fetching problem results:", err)
				utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch problem results", nil)
				return
			}
		}

		// Count total solvers per problem
		type solverCount struct {
			ProblemId   string
			SolvedCount int
		}
		var solverCounts []solverCount
		if err := h.db.Model(&models.ContestProblemResult{}).
			Select("problem_id, COUNT(*) as solved_count").
			Where("contest_id = ? AND is_solved = ?", contestId, true).
			Group("problem_id").Scan(&solverCounts); err != nil {
			log.Println("Error counting solvers:", err)
		}

		solverMap := make(map[string]int)
		for _, sc := range solverCounts {
			solverMap[sc.ProblemId] = sc.SolvedCount
		}

		// Build user result map
		userResultMap := make(map[string]*models.ContestProblemResult)
		for i := range problemResults {
			userResultMap[problemResults[i].ProblemId] = &problemResults[i]
		}

		for _, cp := range contestProblems {
			prob := Problem{
				Id:    cp.ProblemID,
				Index: cp.Index,
			}

			if ur, exists := userResultMap[cp.ProblemID]; exists {
				prob.Solved = ur.IsSolved
				prob.Attempted = ur.WrongAttempts > 0
			}

			if sc, exists := solverMap[cp.ProblemID]; exists {
				prob.TotalSolvers = sc
			}

			problems = append(problems, prob)
		}
	}

	if problems == nil {
		problems = []Problem{}
	}

	response := struct {
		Contest  Contest   `json:"contest"`
		Problems []Problem `json:"problems"`
	}{
		Contest:  contest,
		Problems: problems,
	}

	utils.SendResponse(w, http.StatusOK, "Contest fetched successfully", response)
}
