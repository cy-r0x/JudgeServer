package models

import (
	"time"
)

type Problem struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title              string    `gorm:"type:varchar(255);not null" json:"title"`
	Slug               string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Statement          string    `gorm:"type:text;not null" json:"statement"`
	InputStatement     string    `gorm:"type:text;not null" json:"inputStatement"`
	OutputStatement    string    `gorm:"type:text;not null" json:"outputStatement"`
	TimeLimit          float64   `gorm:"not null" json:"timeLimit"`
	MemoryLimit        float64   `gorm:"not null" json:"memoryLimit"`
	CheckerType        string    `gorm:"type:varchar(10);not null" json:"checkerType"`
	CheckerStrictSpace bool      `gorm:"not null" json:"checkerStrictSpace"`
	CheckerPrecision   *string   `gorm:"type:varchar(10)" json:"checkerPrecision"`
	CreatedByID        *uint     `gorm:"column:created_by;index:idx_problems_created_by" json:"createdById"`
	CreatedBy          *User     `gorm:"foreignKey:CreatedByID;constraint:OnDelete:SET NULL;" json:"createdBy,omitempty"`
	CreatedAt          time.Time `gorm:"type:timestamptz;default:now()" json:"createdAt"`
}

type ContestProblem struct {
	ContestID uint     `gorm:"primaryKey;autoIncrement:false;index:idx_contest_problems_contest" json:"contestId"`
	ProblemID uint     `gorm:"primaryKey;autoIncrement:false" json:"problemId"`
	Index     int      `gorm:"not null;index:idx_contest_problems_contest" json:"index"`
	Contest   *Contest `gorm:"foreignKey:ContestID;constraint:OnDelete:CASCADE;" json:"contest,omitempty"`
	Problem   *Problem `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE;" json:"problem,omitempty"`
}
