package compilerun

import (
	"github.com/judgenot0/judge-backend/config"
	"gorm.io/gorm"
)

type Handler struct {
	db     *gorm.DB
	config *config.Config
}

type Testcase struct {
	Input          string `json:"input" gorm:"column:input"`
	ExpectedOutput string `json:"expected_output" gorm:"column:expected_output"`
}

type UserSubmission struct {
	ProblemId  string `json:"problem_id"`
	ContestId  string `json:"contest_id"`
	Language   string `json:"language"`
	SourceCode string `json:"source_code"`
}

type Problem struct {
	Language           string     `json:"language" gorm:"-"`
	SourceCode         string     `json:"source_code" gorm:"-"`
	TimeLimit          float32    `json:"time_limit" gorm:"column:time_limit"`
	MemoryLimit        float32    `json:"memory_limit" gorm:"column:memory_limit"`
	Testcases          []Testcase `json:"testcases" gorm:"-"`
	CheckerType        string     `json:"checker_type" gorm:"column:checker_type"`
	CheckerStrictSpace bool       `json:"checker_strict_space" gorm:"column:checker_strict_space"`
	CheckerPrecision   *string    `json:"checker_precision" gorm:"column:checker_precision"`
}
