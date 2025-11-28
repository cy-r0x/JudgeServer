package standings

import (
	"database/sql"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
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
	FullName     string          `json:"full_name"`
	Clan         *string         `json:"clan"`
	TotalPenalty int             `json:"total_penalty"`
	SolvedCount  int             `json:"solved_count"`
	Problems     []ProblemStatus `json:"problems"`
	LastSolvedAt *time.Time      `json:"last_solved_at,omitempty"`
}

type ProblemSolveStatus struct {
	Solved    int `json:"solved"`
	Attempted int `json:"attempted"`
}

type StandingsResponse struct {
	ContestId          int64                      `json:"contest_id"`
	ContestTitle       string                     `json:"contest_title"`
	TotalProblemCount  int                        `json:"total_problem_count"`
	Standings          []UserStanding             `json:"standings"`
	StartTime          time.Time                  `json:"start_time"`
	DurationSeconds    int64                      `json:"duration_seconds"`
	ProblemSolveStatus map[int]ProblemSolveStatus `json:"problem_solve_status"`
	TotalItem          int                        `json:"total_item"`
	TotalPages         int                        `json:"total_page"`
	Limit              int                        `json:"limit"`
	Page               int                        `json:"page"`
}

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

type Handler struct {
	db             *sqlx.DB
	mu             sync.RWMutex
	Last_standings map[int64]struct {
		timestamp *time.Time
		standings *StandingsResponse
	}
}
