package standings

import (
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
	Name         string          `json:"name"`
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

type Handler struct {
	db             *gorm.DB
	mu             sync.RWMutex
	Last_standings map[string]struct {
		timestamp *time.Time
		standings *StandingsResponse
	}
}