package models

import (
	"time"
)

type Testcase struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProblemID      uint      `gorm:"not null;index:idx_testcases_problem" json:"problemId"`
	Input          string    `gorm:"type:text;not null" json:"input"`
	ExpectedOutput string    `gorm:"type:text;not null" json:"expectedOutput"`
	IsSample       bool      `gorm:"default:false" json:"isSample"`
	CreatedAt      time.Time `gorm:"type:timestamptz;default:now()" json:"createdAt"`

	Problem Problem `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"problem,omitempty"`
}
