package models

import "time"

type ContestProblemResult struct {
	Id                   string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ContestId            string     `gorm:"type:uuid;not null;index:idx_contest_problem_results_contest" json:"contestId"`
	UserId               string     `gorm:"type:uuid;not null;index:idx_contest_problem_results_user" json:"userId"`
	ProblemId            string     `gorm:"type:uuid;not null;index:idx_contest_problem_results_problem" json:"problemId"`
	IsSolved             bool       `gorm:"not null" json:"isSolved"`
	WrongAttempts        int        `gorm:"not null" json:"wrongAttempts"`
	AcceptedSubmissionId *int64     `gorm:"index:idx_contest_problem_results_accepted_submission" json:"acceptedSubmissionId"`
	SolvedAt             *time.Time `gorm:"type:timestamptz" json:"solvedAt"`
	Penalty              int        `gorm:"not null" json:"penalty"`
	IsFirstBlood         bool       `gorm:"not null" json:"isFirstBlood"`

	Contest            Contest     `gorm:"foreignKey:ContestId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"contest"`
	User               User        `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Problem            Problem     `gorm:"foreignKey:ProblemId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"problem"`
	AcceptedSubmission *Submission `gorm:"foreignKey:AcceptedSubmissionId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"acceptedSubmission,omitempty"`
}
