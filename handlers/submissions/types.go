package submissions

import (
	"time"
)

type Submission struct {
	Id            int64     `json:"id" db:"id"`
	UserId        string    `json:"user_id,omitempty" db:"user_id"`
	Username      string    `json:"username" db:"username"`
	Clan          *string   `json:"clan,omitempty" db:"clan"`
	FullName      string    `json:"full_name,omitempty" db:"full_name"`
	RoomNo        *string   `json:"room_no,omitempty" db:"room_no"`
	PcNo          *string   `json:"pc_no,omitempty" db:"pc_no"`
	ProblemId     string    `json:"problem_id" db:"problem_id"`
	ProblemIndex  int       `json:"problem_index" db:"problem_index"`
	ContestId     string    `json:"contest_id,omitempty" db:"contest_id"`
	Language      string    `json:"language" db:"language"`
	SourceCode    string    `json:"source_code,omitempty" db:"source_code"`
	Verdict       string    `json:"verdict" db:"verdict"`
	ExecutionTime *float32  `json:"execution_time" db:"execution_time"`
	MemoryUsed    *float32  `json:"memory_used" db:"memory_used"`
	SubmittedAt   time.Time `json:"submitted_at" db:"submitted_at"`
}

type UserSubmission struct {
	ProblemId  string `json:"problem_id"`
	ContestId  string `json:"contest_id"`
	Language   string `json:"language"`
	SourceCode string `json:"source_code"`
}

type Testcase struct {
	Input          string `json:"input" db:"input"`
	ExpectedOutput string `json:"expected_output" db:"expected_output"`
}

type Problem struct {
	SubmissionId       int64      `gorm:"-" json:"submission_id"`
	Language           string     `gorm:"-" json:"language"`
	SourceCode         string     `gorm:"-" json:"source_code"`
	TimeLimit          float32    `json:"time_limit" db:"time_limit"`
	MemoryLimit        float32    `json:"memory_limit" db:"memory_limit"`
	Testcases          []Testcase `gorm:"-" json:"testcases"`
	CheckerStrictSpace bool       `json:"checker_strict_space" db:"checker_strict_space"`
	CheckerPrecision   *string    `json:"checker_precision" db:"checker_precision"`
}
