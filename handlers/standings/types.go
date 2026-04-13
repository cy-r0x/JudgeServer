package standings

import (
	"database/sql"
	"sync"
	"time"

	"gorm.io/gorm"
)

type ProblemStatus struct {
	Solved        bool       `json:"solved"`
	FirstSolvedAt *time.Time `json:"first_solved_at,omitempty"`
	Attempts      int        `json:"attempts"`
	Penalty       int        `json:"penalty"`
	FirstBlood    bool       `json:"first_blood"`
}

type UserStanding struct {
	UserId       string          `json:"user_id"`
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
	ContestId          string                     `json:"contest_id"`
	ContestTitle       string                     `json:"contest_title"`
	ProblemMapping     map[int]string             `json:"problem_mapping"`
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
	ProblemId string `gorm:"column:problem_id"`
	Index     int    `gorm:"column:index"`
	Title     string `gorm:"column:title"`
}

type ContestInfo struct {
	Title           string    `gorm:"column:title"`
	StartTime       time.Time `gorm:"column:start_time"`
	DurationSeconds int64     `gorm:"column:duration_seconds"`
}

type userStandingRow struct {
	UserId       string       `gorm:"column:user_id"`
	Username     string       `gorm:"column:username"`
	FullName     string       `gorm:"column:full_name"`
	Clan         *string      `gorm:"column:clan"`
	SolvedCount  int          `gorm:"column:solved_count"`
	TotalPenalty int          `gorm:"column:penalty"`
	LastSolvedAt sql.NullTime `gorm:"column:last_solved_at"`
}

type userProblemRow struct {
	UserId       string       `gorm:"column:user_id"`
	ProblemId    string       `gorm:"column:problem_id"`
	ProblemIndex int          `gorm:"column:problem_index"`
	IsSolved     bool         `gorm:"column:is_solved"`
	SolvedAt     sql.NullTime `gorm:"column:solved_at"`
	Penalty      int          `gorm:"column:penalty"`
	AttemptCount int          `gorm:"column:attempt_count"`
	FirstBlood   bool         `gorm:"column:first_blood"`
}

type problemStatsRow struct {
	ProblemIndex   int `gorm:"column:problem_index"`
	SolvedCount    int `gorm:"column:solved_count"`
	AttemptedUsers int `gorm:"column:attempted_users"`
}

type Handler struct {
	db             *gorm.DB
	mu             sync.RWMutex
	Last_standings map[string]struct {
		timestamp *time.Time
		standings *StandingsResponse
	}
}
