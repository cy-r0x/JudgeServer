package contest

import (
	"time"

	"github.com/judgenot0/judge-backend/config"
	"gorm.io/gorm"
)

type Contest struct {
	Id              string    `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	UserPrefix      string    `json:"userPrefix" db:"user_prefix"`
	Description     string    `json:"description" db:"description"`
	StartTime       time.Time `json:"startTime" db:"start_time"`
	EndTime         time.Time `json:"endTime" db:"end_time"`
	DurationSeconds int64     `json:"durationSeconds" db:"duration_seconds"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

type ContestProblem struct {
	ContestId string `json:"contestId" db:"contest_id"`
	ProblemId string `json:"problemId" db:"problem_id"`
	Index     int    `json:"index" db:"index"`
}

type Handler struct {
	db     *gorm.DB
	config *config.Config
}
