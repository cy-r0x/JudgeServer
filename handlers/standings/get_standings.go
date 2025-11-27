package standings

import (
	"database/sql"
	"fmt"
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

	// Get all problems in the contest
	type ContestProblem struct {
		ProblemId int64 `db:"problem_id"`
		Index     int   `db:"index"`
	}

	var contestProblems []ContestProblem
	err = h.db.Select(&contestProblems, `
		SELECT problem_id, index 
		FROM contest_problems 
		WHERE contest_id = $1 
		ORDER BY index ASC
	`, contestId)
	if err != nil {
		log.Println("Error fetching contest problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contest problems")
		return
	}

	// Fetch contest title
	var contestTitle string
	if err := h.db.Get(&contestTitle, `SELECT title FROM contests WHERE id = $1`, contestId); err != nil {
		log.Println("Error fetching contest title:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contest title")
		return
	}

	type userStandingRow struct {
		UserId       int64         `db:"user_id"`
		Username     string        `db:"username"`
		FullName     string        `db:"full_name"`
		Clan         *string       `db:"clan"`
		SolvedCount  sql.NullInt64 `db:"solved_count"`
		TotalPenalty sql.NullInt64 `db:"penalty"`
	}

	var userStandings []userStandingRow
	err = h.db.Select(&userStandings, `
		WITH wrong_counts AS (
			SELECT 
				user_id,
				COUNT(*) as wrong_count
			FROM submissions
			WHERE contest_id = $1 
			AND LOWER(verdict) IN ('wa','tle','re','mle')
			GROUP BY user_id
		)
		SELECT 
			u.id AS user_id,
			u.username,
			u.full_name,
			u.clan,
			cs.solved_count,
			COALESCE(cs.penalty, COALESCE(wc.wrong_count, 0)) AS penalty
		FROM (
			SELECT DISTINCT user_id 
			FROM submissions 
			WHERE contest_id = $1
		) participants
		JOIN users u ON u.id = participants.user_id
		LEFT JOIN contest_standings cs ON cs.contest_id = $1 AND cs.user_id = participants.user_id
		LEFT JOIN wrong_counts wc ON wc.user_id = participants.user_id
		ORDER BY COALESCE(cs.solved_count, 0) DESC, COALESCE(cs.penalty, COALESCE(wc.wrong_count, 0)) ASC, u.id ASC
	`, contestId)
	if err != nil {
		log.Println("Error fetching user standings:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch standings")
		return
	}

	type problemSolveRow struct {
		UserId       int64         `db:"user_id"`
		ProblemId    int64         `db:"problem_id"`
		ProblemIndex int           `db:"problem_index"`
		SolvedAt     sql.NullTime  `db:"solved_at"`
		Penalty      sql.NullInt64 `db:"penalty"`
		FirstBlood   sql.NullBool  `db:"first_blood"`
		AttemptCount sql.NullInt64 `db:"attempt_count"`
	}

	var problemSolves []problemSolveRow
	err = h.db.Select(&problemSolves, `
		SELECT 
			cs.user_id,
			cs.problem_id,
			cp.index AS problem_index,
			cs.solved_at,
			cs.penalty,
			cs.first_blood,
			cs.attempt_count
		FROM contest_solves cs
		JOIN contest_problems cp ON cp.problem_id = cs.problem_id AND cp.contest_id = cs.contest_id
		WHERE cs.contest_id = $1
		ORDER BY cs.user_id, cp.index
	`, contestId)
	if err != nil {
		log.Println("Error fetching problem solves:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch problem solves")
		return
	}

	problemSolvesMap := make(map[int64][]problemSolveRow)
	for _, ps := range problemSolves {
		problemSolvesMap[ps.UserId] = append(problemSolvesMap[ps.UserId], ps)
	}

	type attemptRow struct {
		UserId    int64 `db:"user_id"`
		ProblemId int64 `db:"problem_id"`
		Attempts  int   `db:"attempts"`
	}

	var attempts []attemptRow
	err = h.db.Select(&attempts, `
		SELECT 
			s.user_id,
			s.problem_id,
			COUNT(*) as attempts
		FROM submissions s
		WHERE s.contest_id = $1
		GROUP BY s.user_id, s.problem_id
	`, contestId)
	if err != nil {
		log.Println("Error fetching attempts:", err)
	}

	attemptsMap := make(map[string]int)
	for _, a := range attempts {
		key := fmt.Sprintf("%d:%d", a.UserId, a.ProblemId)
		attemptsMap[key] = a.Attempts
	}

	standings := make([]UserStanding, 0, len(userStandings))
	for _, us := range userStandings {
		solvedCount := 0
		totalPenalty := 0
		if us.SolvedCount.Valid {
			solvedCount = int(us.SolvedCount.Int64)
		}
		if us.TotalPenalty.Valid {
			totalPenalty = int(us.TotalPenalty.Int64)
		}

		userStanding := UserStanding{
			UserId:       us.UserId,
			Username:     us.Username,
			FullName:     us.FullName,
			Clan:         us.Clan,
			SolvedCount:  solvedCount,
			TotalPenalty: totalPenalty,
			Problems:     make([]ProblemStatus, 0, len(contestProblems)),
		}

		problemSolvesForUser := problemSolvesMap[us.UserId]
		problemSolvesMapByProblem := make(map[int64]problemSolveRow)
		for _, ps := range problemSolvesForUser {
			problemSolvesMapByProblem[ps.ProblemId] = ps
			if ps.SolvedAt.Valid {
				solvedAt := ps.SolvedAt.Time
				if userStanding.LastSolvedAt == nil || solvedAt.After(*userStanding.LastSolvedAt) {
					userStanding.LastSolvedAt = &solvedAt
				}
			}
		}

		for _, cp := range contestProblems {
			ps, hasSolve := problemSolvesMapByProblem[cp.ProblemId]
			attemptKey := fmt.Sprintf("%d:%d", us.UserId, cp.ProblemId)
			attemptCount := attemptsMap[attemptKey]

			problemStatus := ProblemStatus{
				ProblemId:    cp.ProblemId,
				ProblemIndex: cp.Index,
				Attempts:     attemptCount,
				Solved:       hasSolve,
			}

			if hasSolve {
				if ps.SolvedAt.Valid {
					solvedAt := ps.SolvedAt.Time
					problemStatus.FirstSolvedAt = &solvedAt
				}
				if ps.Penalty.Valid {
					problemStatus.Penalty = int(ps.Penalty.Int64)
				}
				if ps.FirstBlood.Valid {
					problemStatus.FirstBlood = ps.FirstBlood.Bool
				}
			}

			userStanding.Problems = append(userStanding.Problems, problemStatus)
		}

		standings = append(standings, userStanding)
	}

	var data struct {
		StartTime time.Time `db:"start_time"`
		Duration  int64     `db:"duration_seconds"`
	}
	err = h.db.Get(&data, `SELECT start_time, duration_seconds FROM contests WHERE id = $1`, contestId)
	if err != nil {
		log.Println("Error fetching contest data:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contest data")
		return
	}

	// Calculate problem statistics (solved and attempted counts per problem index)
	type problemStats struct {
		Solved    int `db:"solved"`
		Attempted int `db:"attempted"`
		Index     int `db:"index"`
	}

	var stats []problemStats
	err = h.db.Select(&stats, `
		SELECT 
			cp.index,
			COUNT(DISTINCT CASE WHEN LOWER(s.verdict) = 'ac' THEN s.user_id END) as solved,
			COUNT(DISTINCT s.user_id) as attempted
		FROM contest_problems cp
		LEFT JOIN submissions s ON s.problem_id = cp.problem_id AND s.contest_id = cp.contest_id
		WHERE cp.contest_id = $1
		GROUP BY cp.index
		ORDER BY cp.index
	`, contestId)
	if err != nil {
		log.Println("Error fetching problem statistics:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch problem statistics")
		return
	}

	report := make(map[int]ProblemReport)
	for _, stat := range stats {
		report[stat.Index] = ProblemReport{
			Solved:    stat.Solved,
			Attempted: stat.Attempted,
		}
	}

	response := StandingsResponse{
		ContestId:         contestId,
		ContestTitle:      contestTitle,
		TotalProblemCount: len(contestProblems),
		Standings:         standings,
		StartTime:         data.StartTime,
		DurationSeconds:   data.Duration,
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
