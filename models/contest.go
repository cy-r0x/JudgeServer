package models

import (
	"time"
)

type Contest struct {
	Id              string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title           string    `gorm:"type:varchar(255);not null" json:"title"`
	UserPrefix      string    `gorm:"type:varchar(50);not null;uniqueIndex:uq_contest_user_prefix" json:"userPrefix"`
	Description     *string   `gorm:"type:text" json:"description"`
	StartTime       time.Time `gorm:"type:timestamptz;index:idx_contests_start_time;not null" json:"startTime"`
	EndTime         time.Time `gorm:"type:timestamptz;index:idx_contests_end_time;not null" json:"endTime"`
	DurationSeconds int64     `gorm:"not null" json:"durationSeconds"`
	CreatedAt       time.Time `gorm:"type:timestamptz;default:now()" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"type:timestamptz;default:now()" json:"updatedAt"`
}
