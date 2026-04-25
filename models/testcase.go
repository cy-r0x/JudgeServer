package models

import (
	"time"
)

type Testcase struct {
	Id             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProblemID      string    `gorm:"type:uuid;not null;index:idx_testcases_problem" json:"problemId"`
	Input          string    `gorm:"type:text;not null" json:"input"`
	ExpectedOutput string    `gorm:"type:text;not null" json:"expectedOutput"`
	IsSample       bool      `gorm:"default:false" json:"isSample"`
	CreatedAt      time.Time `gorm:"type:timestamptz;default:now()" json:"createdAt"`

	Problem Problem `gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE" json:"problem"`
}
