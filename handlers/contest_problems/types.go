package contest_problems

import (
	"gorm.io/gorm"
)

type ContestProblem struct {
	ContestId     string `json:"contest_id" gorm:"column:contest_id"`
	ProblemId     string `json:"problem_id" gorm:"column:problem_id"`
	Index         int    `json:"index" gorm:"column:index"`
	ProblemName   string `json:"problem_name,omitempty" gorm:"column:problem_name"`
	ProblemAuthor string `json:"problem_author,omitempty" gorm:"column:problem_author"`
}

type Handler struct {
	db *gorm.DB
}
