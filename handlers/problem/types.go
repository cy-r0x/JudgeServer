package problem

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Testcase struct {
	Id             int64     `json:"id" db:"id"`
	ProblemId      int64     `json:"problem_id" db:"problem_id"`
	Input          string    `json:"input" db:"input"`
	ExpectedOutput string    `json:"expected_output" db:"expected_output"`
	IsSample       bool      `json:"is_sample" db:"is_sample"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type Problem struct {
	Id                 int64      `json:"id" db:"id"`
	Title              string     `json:"title" db:"title"`
	Slug               string     `json:"slug" db:"slug"`
	Statement          string     `json:"statement" db:"statement"`
	InputStatement     string     `json:"input_statement" db:"input_statement"`
	OutputStatement    string     `json:"output_statement" db:"output_statement"`
	TimeLimit          float32    `json:"time_limit" db:"time_limit"`
	MemoryLimit        float32    `json:"memory_limit" db:"memory_limit"`
	Testcases          []Testcase `json:"test_cases"`
	CheckerType        string     `json:"checker_type" db:"checker_type"`
	CheckerStrictSpace bool       `json:"checker_strict_space" db:"checker_strict_space"`
	CheckerPrecision   *string    `json:"checker_precision" db:"checker_precision"`
	StartTime          *time.Time `json:"start_time,omitempty" db:"start_time"`
	DurationSeconds    *int64     `json:"duration_seconds,omitempty" db:"duration_seconds"`
	CreatedBy          int64      `json:"created_by" db:"created_by"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
}

type Handler struct {
	db *sqlx.DB
}
