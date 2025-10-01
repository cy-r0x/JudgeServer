package contest

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Contest struct {
	Id              int64     `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	StartTime       time.Time `json:"start_time" db:"start_time"`
	DurationSeconds int64     `json:"duration_seconds" db:"duration_seconds"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type ContestProblem struct {
	ContestId int64 `json:"contest_id" db:"contest_id"`
	ProblemId int64 `json:"problem_id" db:"problem_id"`
	Index     int   `json:"index" db:"index"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}
