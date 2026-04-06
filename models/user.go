package models

import (
	"time"
)

type User struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName       string    `gorm:"type:varchar(100);not null" json:"fullName"`
	Username       string    `gorm:"type:varchar(50);uniqueIndex:unique_username_per_contest,priority:1;not null" json:"username"`
	Password       string    `gorm:"type:text;not null" json:"-"`
	Role           string    `gorm:"type:varchar(20);not null;default:'user';check:role IN ('user', 'setter', 'admin')" json:"role"`
	AllowedContest *uint     `gorm:"uniqueIndex:unique_username_per_contest,priority:2" json:"allowedContest"`
	Clan           *string   `gorm:"type:varchar(255)" json:"clan"`
	RoomNo         *string   `gorm:"type:varchar(50)" json:"roomNo"`
	PcNo           *string   `gorm:"type:varchar(50)" json:"pcNo"`
	CreatedAt      time.Time `gorm:"type:timestamptz;default:now()" json:"createdAt"`
}
