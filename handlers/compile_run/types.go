package compilerun

import (
	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/config"
)

type Handler struct {
	db     *sqlx.DB
	config *config.Config
}

type Testcase struct {
	Input          string `json:"input" db:"input"`
	ExpectedOutput string `json:"expected_output" db:"expected_output"`
}

type UserSubmission struct {
	ProblemId  int64  `json:"problem_id"`
	ContestId  int64  `json:"contest_id"`
	Language   string `json:"language"`
	SourceCode string `json:"source_code"`
}

type Problem struct {
	Language           string     `json:"language"`
	SourceCode         string     `json:"source_code"`
	TimeLimit          float32    `json:"time_limit" db:"time_limit"`
	MemoryLimit        float32    `json:"memory_limit" db:"memory_limit"`
	Testcases          []Testcase `json:"testcases"`
	CheckerType        string     `json:"checker_type" db:"checker_type"`
	CheckerStrictSpace bool       `json:"checker_strict_space" db:"checker_strict_space"`
	CheckerPrecision   *string    `json:"checker_precision" db:"checker_precision"`
}
