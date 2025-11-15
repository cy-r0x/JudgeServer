package submissions

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/config"
)

type Submission struct {
	Id            int64     `json:"id" db:"id"`
	UserId        int64     `json:"user_id,omitempty" db:"user_id"`
	Username      string    `json:"username" db:"username"`
	Clan          *string   `json:"clan,omitempty" db:"clan"`
	FullName      string    `json:"full_name,omitempty" db:"full_name"`
	RoomNo        *string   `json:"room_no,omitempty" db:"room_no"`
	PcNo          *string   `json:"pc_no,omitempty" db:"pc_no"`
	ProblemId     int64     `json:"problem_id" db:"problem_id"`
	ContestId     int64     `json:"contest_id,omitempty" db:"contest_id"`
	Language      string    `json:"language" db:"language"`
	SourceCode    string    `json:"source_code,omitempty" db:"source_code"`
	Verdict       string    `json:"verdict" db:"verdict"`
	ExecutionTime *float32  `json:"execution_time" db:"execution_time"`
	MemoryUsed    *float32  `json:"memory_used" db:"memory_used"`
	SubmittedAt   time.Time `json:"submitted_at" db:"submitted_at"`
}

type UserSubmission struct {
	ProblemId  int64  `json:"problem_id"`
	ContestId  int64  `json:"contest_id"`
	Language   string `json:"language"`
	SourceCode string `json:"source_code"`
}

type Testcase struct {
	Input          string `json:"input" db:"input"`
	ExpectedOutput string `json:"expected_output" db:"expected_output"`
}

type QueueSubmission struct {
	SubmissionId int64      `json:"submission_id"`
	ProblemId    int64      `json:"problem_id"`
	Language     string     `json:"language"`
	SourceCode   string     `json:"source_code"`
	Testcases    []Testcase `json:"testcases"`
	Timelimit    float32    `json:"time_limit"`
	MemoryLimit  float32    `json:"memory_limit"`
}

type Handler struct {
	db     *sqlx.DB
	config *config.Config
}

func NewHandler(db *sqlx.DB, config *config.Config) *Handler {
	return &Handler{
		db:     db,
		config: config,
	}
}
