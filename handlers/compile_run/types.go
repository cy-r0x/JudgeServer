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
	ExpectedOutput string `json:"expectedOutput" gorm:"column:expected_output"`
}

type UserSubmission struct {
	ProblemId  string `json:"problemId"`
	ContestId  string `json:"contestId"`
	Language   string `json:"language"`
	SourceCode string `json:"sourceCode"`
}

type Problem struct {
	Language           string     `json:"language" gorm:"-"`
	SourceCode         string     `json:"source_code" gorm:"-"`
	TimeLimit          float32    `json:"timeLimit" gorm:"column:time_limit"`
	MemoryLimit        float32    `json:"memoryLimit" gorm:"column:memory_limit"`
	Testcases          []Testcase `json:"testcases" gorm:"-"`
	CheckerType        string     `json:"checkerType" gorm:"column:checker_type"`
	CheckerStrictSpace bool       `json:"checkerStrictSpace" gorm:"column:checker_strict_space"`
	CheckerPrecision   *string    `json:"checkerPrecision" gorm:"column:checker_precision"`
}
