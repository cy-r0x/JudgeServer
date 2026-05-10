package contest_problems

import (
	"gorm.io/gorm"
)

type ContestProblem struct {
	ContestId     string `json:"contestId" gorm:"column:contest_id"`
	ProblemId     string `json:"problemId" gorm:"column:problem_id"`
	Index         int    `json:"index" gorm:"column:index"`
	ProblemName   string `json:"problemName,omitempty" gorm:"column:problem_name"`
	ProblemAuthor string `json:"problemAuthor,omitempty" gorm:"column:problem_author"`
}

type Handler struct {
	db *gorm.DB
}
