package standings

import (
	"log"
	"net/http"
	"sync"

	"github.com/judgenot0/judge-backend/utils"
)

type ExportedUserStanding struct {
	FullName   string  `json:"full_name"`
	Clan       *string `json:"clan"`
	SolveCount int     `json:"solve_count"`
	Solved     []int   `json:"solved"`
}

func (h *Handler) ExportStandings(w http.ResponseWriter, r *http.Request) {
	contestIdStr := r.PathValue("contestId")
	if contestIdStr == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID is required")
		return
	}
	contestId := contestIdStr

	// Fetch all data in parallel using goroutines
	var wg sync.WaitGroup
	var contestProblems []ContestProblem
	var userStandings []userStandingRow
	var userProblems []userProblemRow
	var errChan = make(chan error, 3)

	// Query 1: Fetch contest problems
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.db.Raw(`
			SELECT cp.problem_id, cp.index, p.title
			FROM contest_problems cp
			JOIN problems p ON p.id = cp.problem_id
			WHERE cp.contest_id = ? 
			ORDER BY cp.index ASC
		`, contestId).Scan(&contestProblems).Error
		if err != nil {
			errChan <- err
		}
	}()

	// Query 2: Fetch user standings
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.db.Raw(`
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
				WHERE contest_id = ?
			) participants
			JOIN users u ON u.id = participants.user_id
			LEFT JOIN contest_standings cs ON cs.contest_id = ? AND cs.user_id = participants.user_id
			ORDER BY COALESCE(cs.solved_count, 0) DESC, (COALESCE(cs.penalty, 0) + COALESCE(cs.wrong_attempts, 0)) ASC, u.id ASC
		`, contestId, contestId).Scan(&userStandings).Error
		if err != nil {
			errChan <- err
		}
	}()

	// Query 3: Fetch user-problem details
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := h.db.Raw(`
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
			WHERE contest_id = ?
			ORDER BY user_id, problem_index
		`, contestId).Scan(&userProblems).Error
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
	userProblemsMap := make(map[string]map[string]userProblemRow)
	for _, up := range userProblems {
		if _, exists := userProblemsMap[up.UserId]; !exists {
			userProblemsMap[up.UserId] = make(map[string]userProblemRow)
		}
		userProblemsMap[up.UserId][up.ProblemId] = up
	}

	// Build exported standings
	exportedStandings := make([]ExportedUserStanding, 0, len(userStandings))
	for _, us := range userStandings {
		solved := make([]int, 0)
		userProblemData := userProblemsMap[us.UserId]

		for _, cp := range contestProblems {
			up, hasData := userProblemData[cp.ProblemId]
			if hasData && up.IsSolved {
				solved = append(solved, cp.Index)
			}
		}

		exportedStanding := ExportedUserStanding{
			FullName:   us.FullName,
			Clan:       us.Clan,
			SolveCount: us.SolvedCount,
			Solved:     solved,
		}

		exportedStandings = append(exportedStandings, exportedStanding)
	}

	utils.SendResponse(w, http.StatusOK, exportedStandings)
}
