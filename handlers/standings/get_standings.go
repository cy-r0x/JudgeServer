package standings

import (
	"log"
	"net/http"
	"strconv"
	"sync"
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

			response.TotalItem = totalStandings
			response.TotalPages = (totalStandings + limit - 1) / limit
			response.Limit = limit
			response.Page = crrPage

			utils.SendResponse(w, http.StatusOK, response)
			return
		}
	}

	// Fetch all data in parallel using goroutines
	var wg sync.WaitGroup
	var contestInfo ContestInfo
	var contestProblems []ContestProblem
	var userStandings []userStandingRow
	var userProblems []userProblemRow
	var stats []problemStatsRow
	var errChan = make(chan error, 5)

	// Query 1: Fetch contest info
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.db.Get(&contestInfo, `SELECT title, start_time, duration_seconds FROM contests WHERE id = $1`, contestId)
		if err != nil {
			errChan <- err
		}
	}()

	// Query 2: Fetch contest problems
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.db.Select(&contestProblems, `
			SELECT problem_id, index, title
			FROM contest_problems cp
			JOIN problems p ON p.id = cp.problem_id
			WHERE contest_id = $1 
			ORDER BY index ASC
		`, contestId)
		if err != nil {
			errChan <- err
		}
	}()

	// Query 3: Fetch user standings
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.db.Select(&userStandings, `
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
			errChan <- err
		}
	}()

	// Query 4: Fetch user-problem details
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.db.Select(&userProblems, `
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
			errChan <- err
		}
	}()

	// Query 5: Fetch problem statistics
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.db.Select(&stats, `
			SELECT 
				problem_index,
				solved_count,
				attempted_users
			FROM contest_problem_stats
			WHERE contest_id = $1
			ORDER BY problem_index
		`, contestId)
		if err != nil {
			errChan <- err
		}
	}()

	// Wait for all queries to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		log.Println("Database query error:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch standings data")
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
	problem_solve_status := make(map[int]ProblemSolveStatus)
	for _, stat := range stats {
		problem_solve_status[stat.ProblemIndex] = ProblemSolveStatus{
			Solved:    stat.SolvedCount,
			Attempted: stat.AttemptedUsers,
		}
	}

	response := StandingsResponse{
		ContestId:          contestId,
		ContestTitle:       contestInfo.Title,
		TotalProblemCount:  len(contestProblems),
		Standings:          standings,
		StartTime:          contestInfo.StartTime,
		DurationSeconds:    contestInfo.DurationSeconds,
		ProblemSolveStatus: problem_solve_status,
		TotalItem:          len(standings),
		TotalPages:         (len(standings) + limit - 1) / limit,
		Limit:              limit,
		Page:               crrPage,
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
