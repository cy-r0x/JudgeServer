package standings

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/utils"
)

type ProblemStatus struct {
	ProblemId     int64      `json:"problem_id"`
	ProblemIndex  int        `json:"problem_index"`
	Solved        bool       `json:"solved"`
	FirstSolvedAt *time.Time `json:"first_solved_at,omitempty"`
	Attempts      int        `json:"attempts"`
	Penalty       int        `json:"penalty"`
	FirstBlood    bool       `json:"first_blood"`
}

type UserStanding struct {
	UserId       int64           `json:"user_id"`
	Username     string          `json:"username"`
	TotalPenalty int             `json:"total_penalty"`
	SolvedCount  int             `json:"solved_count"`
	Problems     []ProblemStatus `json:"problems"`
	LastSolvedAt *time.Time      `json:"last_solved_at,omitempty"`
}

type StandingsResponse struct {
	ContestId         int64          `json:"contest_id"`
	TotalProblemCount int            `json:"total_problem_count"`
	Standings         []UserStanding `json:"standings"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetStandings(w http.ResponseWriter, r *http.Request) {
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

	// Get all users who participated (made at least one submission)
	type User struct {
		UserId   int64  `db:"user_id"`
		Username string `db:"username"`
	}

	var users []User
	err = h.db.Select(&users, `
		SELECT DISTINCT u.id as user_id, u.username 
		FROM users u
		INNER JOIN submissions s ON u.id = s.user_id
		WHERE s.contest_id = $1
		ORDER BY u.id
	`, contestId)
	if err != nil {
		log.Println("Error fetching users:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	// Get all solves
	type SolveData struct {
		UserId       int64     `db:"user_id"`
		ProblemId    int64     `db:"problem_id"`
		SolvedAt     time.Time `db:"solved_at"`
		Penalty      int       `db:"penalty"`
		AttemptCount int       `db:"attempt_count"`
		FirstBlood   bool      `db:"first_blood"`
	}

	var solves []SolveData
	err = h.db.Select(&solves, `
		SELECT user_id, problem_id, solved_at, penalty, attempt_count, first_blood
		FROM contest_solves
		WHERE contest_id = $1
	`, contestId)
	if err != nil {
		log.Println("Error fetching solves:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch solves")
		return
	}

	// Build a map of user -> problem -> solve data
	userSolveMap := make(map[int64]map[int64]*SolveData)
	for i := range solves {
		if userSolveMap[solves[i].UserId] == nil {
			userSolveMap[solves[i].UserId] = make(map[int64]*SolveData)
		}
		userSolveMap[solves[i].UserId][solves[i].ProblemId] = &solves[i]
	}

	// Get attempt counts for all users and problems (including unsolved)
	type AttemptData struct {
		UserId    int64 `db:"user_id"`
		ProblemId int64 `db:"problem_id"`
		Attempts  int   `db:"attempts"`
	}

	var attempts []AttemptData
	err = h.db.Select(&attempts, `
		SELECT s.user_id, s.problem_id, COUNT(*) as attempts
		FROM submissions s
		LEFT JOIN contest_solves cs ON cs.contest_id = s.contest_id 
			AND cs.user_id = s.user_id 
			AND cs.problem_id = s.problem_id
		WHERE s.contest_id = $1
		AND (cs.solved_at IS NULL OR s.submitted_at <= cs.solved_at)
		GROUP BY s.user_id, s.problem_id
	`, contestId)
	if err != nil {
		log.Println("Error fetching attempts:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch attempts")
		return
	}

	// Build a map of user -> problem -> attempt count
	userAttemptMap := make(map[int64]map[int64]int)
	for _, attempt := range attempts {
		if userAttemptMap[attempt.UserId] == nil {
			userAttemptMap[attempt.UserId] = make(map[int64]int)
		}
		userAttemptMap[attempt.UserId][attempt.ProblemId] = attempt.Attempts
	}

	// Build the standings
	standings := make([]UserStanding, 0, len(users))

	for _, user := range users {
		problems := make([]ProblemStatus, 0, len(contestProblems))
		totalPenalty := 0
		solvedCount := 0
		var lastSolvedAt *time.Time

		for _, cp := range contestProblems {
			problemStatus := ProblemStatus{
				ProblemId:    cp.ProblemId,
				ProblemIndex: cp.Index,
				Solved:       false,
				Attempts:     0,
				Penalty:      0,
				FirstBlood:   false,
			}

			// Get attempt count
			if userAttemptMap[user.UserId] != nil {
				problemStatus.Attempts = userAttemptMap[user.UserId][cp.ProblemId]
			}

			// Check if solved
			if userSolveMap[user.UserId] != nil {
				if solveData, exists := userSolveMap[user.UserId][cp.ProblemId]; exists {
					problemStatus.Solved = true
					problemStatus.FirstSolvedAt = &solveData.SolvedAt
					problemStatus.Penalty = solveData.Penalty
					problemStatus.FirstBlood = solveData.FirstBlood

					totalPenalty += solveData.Penalty
					solvedCount++

					if lastSolvedAt == nil || solveData.SolvedAt.After(*lastSolvedAt) {
						lastSolvedAt = &solveData.SolvedAt
					}
				}
			}

			problems = append(problems, problemStatus)
		}

		standings = append(standings, UserStanding{
			UserId:       user.UserId,
			Username:     user.Username,
			TotalPenalty: totalPenalty,
			SolvedCount:  solvedCount,
			Problems:     problems,
			LastSolvedAt: lastSolvedAt,
		})
	}

	// Sort standings: more solves -> lower penalty -> earlier last solve
	sort.SliceStable(standings, func(i, j int) bool {
		if standings[i].SolvedCount != standings[j].SolvedCount {
			return standings[i].SolvedCount > standings[j].SolvedCount
		}
		if standings[i].TotalPenalty != standings[j].TotalPenalty {
			return standings[i].TotalPenalty < standings[j].TotalPenalty
		}
		li, lj := standings[i].LastSolvedAt, standings[j].LastSolvedAt
		switch {
		case li == nil && lj == nil:
			return standings[i].UserId < standings[j].UserId
		case li == nil:
			return false
		case lj == nil:
			return true
		default:
			return li.Before(*lj)
		}
	})

	response := StandingsResponse{
		ContestId:         contestId,
		TotalProblemCount: len(contestProblems),
		Standings:         standings,
	}

	utils.SendResponse(w, http.StatusOK, response)
}
