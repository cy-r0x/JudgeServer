package problem

import (
	"time"

	"gorm.io/gorm"
)

type Testcase struct {
	Id             string    `json:"id" gorm:"column:id"`
	ProblemId      string    `json:"problemId" gorm:"column:problem_id"`
	Input          string    `json:"input" gorm:"column:input"`
	ExpectedOutput string    `json:"expectedOutput" gorm:"column:expected_output"`
	IsSample       bool      `json:"isSample" gorm:"column:is_sample"`
	CreatedAt      time.Time `json:"createdAt" gorm:"column:created_at"`
}

type Problem struct {
	Id                 string              `json:"id" gorm:"column:id"`
	Title              string              `json:"title" gorm:"column:title"`
	Statement          string              `json:"statement" gorm:"column:statement"`
	InputStatement     string              `json:"inputStatement" gorm:"column:input_statement"`
	OutputStatement    string              `json:"outputStatement" gorm:"column:output_statement"`
	TimeLimit          float32             `json:"timeLimit" gorm:"column:time_limit"`
	MemoryLimit        float32             `json:"memoryLimit" gorm:"column:memory_limit"`
	CheckerType        string              `json:"checkerType" gorm:"column:checker_type"`
	CheckerStrictSpace bool                `json:"checkerStrictSpace" gorm:"column:checker_strict_space"`
	CheckerPrecision   *string             `json:"checkerPrecision" gorm:"column:checker_precision"`
	StartTime          *time.Time          `json:"startTime,omitempty" gorm:"column:start_time"`
	DurationSeconds    *int64              `json:"durationSeconds,omitempty" gorm:"column:duration_seconds"`
	Author             string              `json:"author" gorm:"column:author"`
	CreatedAt          time.Time           `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt          time.Time           `json:"updatedAt" gorm:"column:updated_at"`
	LastSubmission     *LastSubmissionData `json:"lastSubmission,omitempty" gorm:"-"`
}

type UpdateProblemPayload struct {
	Id                 string  `json:"id"`
	Title              string  `json:"title"`
	Statement          string  `json:"statement"`
	InputStatement     string  `json:"inputStatement"`
	OutputStatement    string  `json:"outputStatement"`
	TimeLimit          float64 `json:"timeLimit"`
	MemoryLimit        float64 `json:"memoryLimit"`
	CheckerType        string  `json:"checkerType"`
	CheckerStrictSpace bool    `json:"checkerStrictSpace"`
	CheckerPrecision   *string `json:"checkerPrecision"`
}

type LastSubmissionData struct {
	SourceCode string `json:"sourceCode" gorm:"column:source_code"`
	Language   string `json:"language" gorm:"column:language"`
}

type Handler struct {
	db *gorm.DB
}
