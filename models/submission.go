package models

import (
	"time"
)

type Submission struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        string    `gorm:"type:uuid;not null;index:idx_submissions_user;index:idx_submissions_contest_user_problem;index:idx_submissions_penalty_lookup;index:idx_submissions_contest_user" json:"userId"`
	Username      string    `gorm:"type:varchar(50);not null" json:"username"`
	ProblemID     string    `gorm:"type:uuid;not null;index:idx_submissions_problem;index:idx_submissions_contest_user_problem;index:idx_submissions_penalty_lookup" json:"problemId"`
	ContestID     *string   `gorm:"type:uuid;index:idx_submissions_contest;index:idx_submissions_contest_user_problem;index:idx_submissions_penalty_lookup;index:idx_submissions_contest_user" json:"contestId"`
	Language      string    `gorm:"type:varchar(30);not null" json:"language"`
	SourceCode    string    `gorm:"type:text;not null" json:"sourceCode"`
	FilePath      *string   `gorm:"type:text" json:"filePath"`
	Verdict       string    `gorm:"type:varchar(30);default:'Pending';index:idx_submissions_verdict;index:idx_submissions_penalty_lookup" json:"verdict"`
	FirstBlood    bool      `gorm:"not null;default:false" json:"firstBlood"`
	ExecutionTime *float64  `json:"executionTime"`
	MemoryUsed    *float64  `json:"memoryUsed"`
	SubmittedAt   time.Time `gorm:"type:timestamptz;default:now();index:idx_submissions_submitted_at,sort:desc;index:idx_submissions_contest_user_problem,sort:desc;index:idx_submissions_penalty_lookup" json:"submittedAt"`

	User    User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Problem Problem  `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"problem,omitempty"`
	Contest *Contest `gorm:"foreignKey:ContestID;constraint:OnDelete:SET NULL" json:"contest,omitempty"`
}
