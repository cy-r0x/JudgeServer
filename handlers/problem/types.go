package problem

import (
	"time"

	"gorm.io/gorm"
)

type Testcase struct {
	Id             string    `json:"id" gorm:"column:id"`
	ProblemId      string    `json:"problem_id" gorm:"column:problem_id"`
	Input          string    `json:"input" gorm:"column:input"`
	ExpectedOutput string    `json:"expected_output" gorm:"column:expected_output"`
	IsSample       bool      `json:"is_sample" gorm:"column:is_sample"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
}

type Problem struct {
	Id                 string              `json:"id" gorm:"column:id"`
	Title              string              `json:"title" gorm:"column:title"`
	Slug               string              `json:"slug" gorm:"column:slug"`
	Statement          string              `json:"statement" gorm:"column:statement"`
	InputStatement     string              `json:"input_statement" gorm:"column:input_statement"`
	OutputStatement    string              `json:"output_statement" gorm:"column:output_statement"`
	TimeLimit          float32             `json:"time_limit" gorm:"column:time_limit"`
	MemoryLimit        float32             `json:"memory_limit" gorm:"column:memory_limit"`
	Testcases          []Testcase          `json:"test_cases" gorm:"-"`
	CheckerType        string              `json:"checker_type" gorm:"column:checker_type"`
	CheckerStrictSpace bool                `json:"checker_strict_space" gorm:"column:checker_strict_space"`
	CheckerPrecision   *string             `json:"checker_precision" gorm:"column:checker_precision"`
	StartTime          *time.Time          `json:"start_time,omitempty" gorm:"column:start_time"`
	DurationSeconds    *int64              `json:"duration_seconds,omitempty" gorm:"column:duration_seconds"`
	CreatedBy          string              `json:"created_by" gorm:"column:created_by"`
	CreatedAt          time.Time           `json:"created_at" gorm:"column:created_at"`
	LastSubmission     *LastSubmissionData `json:"last_submission,omitempty" gorm:"-"`
}

type LastSubmissionData struct {
	SourceCode string `json:"source_code" gorm:"column:source_code"`
	Language   string `json:"language" gorm:"column:language"`
}

type Handler struct {
	db *gorm.DB
}
