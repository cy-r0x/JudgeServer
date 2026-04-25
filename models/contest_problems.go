package models

type ContestProblem struct {
	Id        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ContestID string `gorm:"type:uuid;not null;index:idx_contest_problems_contest;uniqueIndex:uq_contest_problem,priority:1" json:"contestId"`
	ProblemID string `gorm:"type:uuid;not null;uniqueIndex:uq_contest_problem,priority:2" json:"problemId"`
	Index     int    `gorm:"not null;index:idx_contest_problems_contest" json:"index"`

	Contest *Contest `gorm:"foreignKey:ContestID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"contest,omitempty"`
	Problem *Problem `gorm:"foreignKey:ProblemID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"problem,omitempty"`
}
