package usercsv

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id               string    `json:"id" db:"id"`
	Name             string    `json:"fullName" db:"full_name"`
	Username         string    `json:"username" db:"username"`
	Password         string    `json:"password,omitempty" db:"password"`
	UnHashedPassword string    `json:"-" db:"-"`
	Role             string    `json:"role,omitempty" db:"role"`
	AdditionalInfo   *string   `json:"additionalInfo,omitempty" db:"additional_info"`
	RoomNo           *string   `json:"roomNo,omitempty" db:"room_no"`
	PcNo             *string   `json:"pcNo,omitempty" db:"pc_no"`
	AllowedContest   *string   `json:"allowedContest,omitempty" db:"allowed_contest"`
	CreatedAt        time.Time `json:"createdAt,omitempty" db:"created_at"`
}

type Handler struct {
	db *gorm.DB
	mu sync.Mutex
}
