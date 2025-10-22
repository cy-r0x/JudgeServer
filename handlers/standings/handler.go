package standings

import (
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

type StandingsResponse struct {
	ContestId         int64          `json:"contest_id"`
	ContestTitle      string         `json:"contest_title"`
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
