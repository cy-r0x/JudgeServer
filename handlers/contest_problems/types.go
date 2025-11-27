package contest_problems

import "github.com/jmoiron/sqlx"

type ContestProblem struct {
	ContestId     int64  `json:"contest_id" db:"contest_id"`
	ProblemId     int64  `json:"problem_id" db:"problem_id"`
	Index         int    `json:"index" db:"index"`
	ProblemName   string `json:"problem_name,omitempty" db:"problem_name"`
	ProblemAuthor string `json:"problem_author,omitempty" db:"problem_author"`
}

type Handler struct {
	db *sqlx.DB
}
