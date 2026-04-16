package models

import (
	"time"
)

type Contest struct {
	ID              string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title           string    `gorm:"type:varchar(255);not null" json:"title"`
	Description     *string   `gorm:"type:text" json:"description"`
	StartTime       time.Time `gorm:"type:timestamptz;index:idx_contests_start_time;not null" json:"startTime"`
	DurationSeconds int64     `gorm:"not null" json:"durationSeconds"`
	CreatedAt       time.Time `gorm:"type:timestamptz;default:now()" json:"createdAt"`
}
