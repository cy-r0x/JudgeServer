package models

import (
	"time"
)

type Role string

const (
	RoleUser   Role = "user"
	RoleSetter Role = "setter"
	RoleAdmin  Role = "admin"
)

type User struct {
	Id             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name           string    `gorm:"type:varchar(100);not null" json:"name"`
	Username       string    `gorm:"type:varchar(50);uniqueIndex:unique_username_per_contest,priority:1;not null" json:"username"`
	Password       string    `gorm:"type:text;not null" json:"password,omitempty"`
	Role           Role      `gorm:"type:varchar(20);not null;default:'user';check:role IN ('user', 'setter', 'admin')" json:"role"`
	AllowedContest *string   `gorm:"type:uuid;uniqueIndex:unique_username_per_contest,priority:2" json:"allowedContest"`
	AdditionalInfo *string   `gorm:"type:varchar(255)" json:"additionalInfo"`
	RoomNo         *string   `gorm:"type:varchar(50)" json:"roomNo"`
	PcNo           *string   `gorm:"type:varchar(50)" json:"pcNo"`
	CreatedAt      time.Time `gorm:"type:timestamptz;default:now()" json:"createdAt"`

	AllowedContestRef *Contest `gorm:"foreignKey:AllowedContest;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"allowedContestRef,omitempty"`
}
