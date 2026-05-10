package submissions

import (
	"time"
)

type SubmissionResponse struct {
	Id         int64     `json:"id"`
	UserId     string    `json:"userId"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	ProblemId  string    `json:"problemId"`
	ContestId  *string   `json:"contestId"`
	Language   string    `json:"language"`
	SourceCode string    `json:"sourceCode"`
	Status     string    `json:"status"`
	ExecTime   *float64  `json:"execTime"`
	ExecMemory *float64  `json:"execMemory"`
	CreatedAt  time.Time `json:"createdAt"`
}

type UserSubmission struct {
	ProblemId  string `json:"problemId"`
	ContestId  string `json:"contestId"`
	Language   string `json:"language"`
	SourceCode string `json:"sourceCode"`
}

type Testcase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
}

type QueueSubmission struct {
	SubmissionId int64
	SourceCode   string
	Testcases    []Testcase
	Language     string
	SubmittedAt  int64
}

type SubmissionListParams struct {
	ContestID      string
	UserID         *string // nil = don't filter by user
	Status         string
	SearchName     string
	SearchUsername string
	Limit          int
	Page           int
}
