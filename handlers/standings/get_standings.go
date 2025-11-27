package standings

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetStandings(w http.ResponseWriter, r *http.Request) {

	const limit = 100
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

	contestIdStr := r.PathValue("contestId")
	if contestIdStr == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID is required")
		return
	}

	contestId, err := strconv.ParseInt(contestIdStr, 10, 64)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid contest ID")
		return
	}

	currentTime := time.Now()

	h.mu.RLock()
	entry, exists := h.Last_standings[contestId]
	h.mu.RUnlock()

	if exists && entry.timestamp != nil {
		if currentTime.Sub(*entry.timestamp) < 15*time.Second {
			response := *entry.standings
			totalStandings := len(response.Standings)

			start := offset
			end := offset + limit
			if start >= totalStandings {
				response.Standings = []UserStanding{}
			} else {
				if end > totalStandings {
					end = totalStandings
				}
				response.Standings = response.Standings[start:end]
			}
			utils.SendResponse(w, http.StatusOK, response)
			return
		}
	}

	// Define types for queries
	type ContestProblem struct {
		ProblemId int64  `db:"problem_id"`
		Index     int    `db:"index"`
		Title     string `db:"title"`
	}

	type ContestInfo struct {
		Title           string    `db:"title"`
		StartTime       time.Time `db:"start_time"`
		DurationSeconds int64     `db:"duration_seconds"`
	}

	type userStandingRow struct {
		UserId       int64        `db:"user_id"`
		Username     string       `db:"username"`
		FullName     string       `db:"full_name"`
		Clan         *string      `db:"clan"`
		SolvedCount  int          `db:"solved_count"`
		TotalPenalty int          `db:"penalty"`
		LastSolvedAt sql.NullTime `db:"last_solved_at"`
	}

	type userProblemRow struct {
		UserId       int64        `db:"user_id"`
		ProblemId    int64        `db:"problem_id"`
		ProblemIndex int          `db:"problem_index"`
		IsSolved     bool         `db:"is_solved"`
		SolvedAt     sql.NullTime `db:"solved_at"`
		Penalty      int          `db:"penalty"`
		AttemptCount int          `db:"attempt_count"`
		FirstBlood   bool         `db:"first_blood"`
	}

	type problemStatsRow struct {
		ProblemIndex   int `db:"problem_index"`
		SolvedCount    int `db:"solved_count"`
		AttemptedUsers int `db:"attempted_users"`
	}

	// Fetch contest info
	var contestInfo ContestInfo
	err = h.db.Get(&contestInfo, `SELECT title, start_time, duration_seconds FROM contests WHERE id = $1`, contestId)
	if err != nil {
		log.Println("Error fetching contest info:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contest info")
		return
	}

	// Fetch contest problems
	var contestProblems []ContestProblem
	err = h.db.Select(&contestProblems, `
		SELECT problem_id, index, title
		FROM contest_problems cp
		JOIN problems p ON p.id = cp.problem_id
		WHERE contest_id = $1 
		ORDER BY index ASC
	`, contestId)
	if err != nil {
		log.Println("Error fetching contest problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contest problems")
		return
	}

	// Fetch user standings
	var userStandings []userStandingRow
	err = h.db.Select(&userStandings, `
		SELECT 
			u.id AS user_id,
			u.username,
			u.full_name,
			u.clan,
			COALESCE(cs.solved_count, 0) AS solved_count,
			COALESCE(cs.penalty, 0) + COALESCE(cs.wrong_attempts, 0) AS penalty,
			cs.last_solved_at
		FROM (
			SELECT DISTINCT user_id 
			FROM submissions 
			WHERE contest_id = $1
		) participants
		JOIN users u ON u.id = participants.user_id
		LEFT JOIN contest_standings cs ON cs.contest_id = $1 AND cs.user_id = participants.user_id
		ORDER BY COALESCE(cs.solved_count, 0) DESC, (COALESCE(cs.penalty, 0) + COALESCE(cs.wrong_attempts, 0)) ASC, u.id ASC
	`, contestId)
	if err != nil {
		log.Println("Error fetching user standings:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch standings")
		return
	}

	// Fetch user-problem details
	var userProblems []userProblemRow
	err = h.db.Select(&userProblems, `
		SELECT 
			user_id,
			problem_id,
			problem_index,
			is_solved,
			solved_at,
			penalty,
			attempt_count,
			first_blood
		FROM contest_user_problems
		WHERE contest_id = $1
		ORDER BY user_id, problem_index
	`, contestId)
	if err != nil {
		log.Println("Error fetching user problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch user problems")
		return
	}

	// Fetch problem statistics
	var stats []problemStatsRow
	err = h.db.Select(&stats, `
		SELECT 
			problem_index,
			solved_count,
			attempted_users
		FROM contest_problem_stats
		WHERE contest_id = $1
		ORDER BY problem_index
	`, contestId)
	if err != nil {
		log.Println("Error fetching problem statistics:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch problem statistics")
		return
	}

	// Build map for efficient lookup
	userProblemsMap := make(map[int64]map[int64]userProblemRow)
	for _, up := range userProblems {
		if _, exists := userProblemsMap[up.UserId]; !exists {
			userProblemsMap[up.UserId] = make(map[int64]userProblemRow)
		}
		userProblemsMap[up.UserId][up.ProblemId] = up
	}

	standings := make([]UserStanding, 0, len(userStandings))
	for _, us := range userStandings {
		userStanding := UserStanding{
			UserId:       us.UserId,
			Username:     us.Username,
			FullName:     us.FullName,
			Clan:         us.Clan,
			SolvedCount:  us.SolvedCount,
			TotalPenalty: us.TotalPenalty,
			Problems:     make([]ProblemStatus, 0, len(contestProblems)),
		}

		if us.LastSolvedAt.Valid {
			userStanding.LastSolvedAt = &us.LastSolvedAt.Time
		}

		userProblemData := userProblemsMap[us.UserId]

		for _, cp := range contestProblems {
			up, hasData := userProblemData[cp.ProblemId]

			problemStatus := ProblemStatus{
				ProblemId:    cp.ProblemId,
				ProblemIndex: cp.Index,
				Attempts:     0,
				Solved:       false,
			}

			if hasData {
				problemStatus.Attempts = up.AttemptCount
				problemStatus.Solved = up.IsSolved
				problemStatus.Penalty = up.Penalty
				problemStatus.FirstBlood = up.FirstBlood

				if up.SolvedAt.Valid {
					problemStatus.FirstSolvedAt = &up.SolvedAt.Time
				}
			}

			userStanding.Problems = append(userStanding.Problems, problemStatus)
		}

		standings = append(standings, userStanding)
	}

	// Build report from stats
	report := make(map[int]ProblemReport)
	for _, stat := range stats {
		report[stat.ProblemIndex] = ProblemReport{
			Solved:    stat.SolvedCount,
			Attempted: stat.AttemptedUsers,
		}
	}

	response := StandingsResponse{
		ContestId:         contestId,
		ContestTitle:      contestInfo.Title,
		TotalProblemCount: len(contestProblems),
		Standings:         standings,
		StartTime:         contestInfo.StartTime,
		DurationSeconds:   contestInfo.DurationSeconds,
		Report:            report,
	}

	h.mu.Lock()
	h.Last_standings[contestId] = struct {
		timestamp *time.Time
		standings *StandingsResponse
	}{
		timestamp: &currentTime,
		standings: &response,
	}
	h.mu.Unlock()

	totalStandings := len(response.Standings)
	start := offset
	end := offset + limit
	if start >= totalStandings {
		response.Standings = []UserStanding{}
	} else {
		if end > totalStandings {
			end = totalStandings
		}
		response.Standings = response.Standings[start:end]
	}

	utils.SendResponse(w, http.StatusOK, response)
}
