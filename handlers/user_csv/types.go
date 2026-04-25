package usercsv

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id               string    `json:"id" db:"id"`
	Name             string    `json:"full_name" db:"full_name"`
	Username         string    `json:"username" db:"username"`
	Password         string    `json:"password,omitempty" db:"password"`
	UnHashedPassword string    `json:"-" db:"-"`
	Role             string    `json:"role,omitempty" db:"role"`
	AdditionalInfo   *string   `json:"additional_info,omitempty" db:"additional_info"`
	RoomNo           *string   `json:"room_no,omitempty" db:"room_no"`
	PcNo             *string   `json:"pc_no,omitempty" db:"pc_no"`
	AllowedContest   *string   `json:"allowed_contest,omitempty" db:"allowed_contest"`
	CreatedAt        time.Time `json:"created_at,omitempty" db:"created_at"`
}

type Handler struct {
	db *gorm.DB
	mu sync.Mutex
}
